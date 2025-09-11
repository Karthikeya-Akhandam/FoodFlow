from flask import Flask, request, jsonify, render_template, session
from flask_sqlalchemy import SQLAlchemy
from flask_cors import CORS
import os
from datetime import datetime
import uuid
import json

app = Flask(__name__)
app.config['SECRET_KEY'] = 'your-secret-key-change-this'
CORS(app)

# Database configuration
DATABASE_URL = "postgresql://postgres:12345678@localhost:5432/foodflow?sslmode=disable"
app.config['SQLALCHEMY_DATABASE_URI'] = DATABASE_URL
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

db = SQLAlchemy(app)

# Models
class ChatbotSession(db.Model):
    __tablename__ = 'chatbot_sessions'
    
    id = db.Column(db.String(36), primary_key=True, default=lambda: str(uuid.uuid4()))
    user_id = db.Column(db.String(36), nullable=True)  # Optional - link to your users table
    session_id = db.Column(db.String(100), unique=True, nullable=False)
    created_at = db.Column(db.DateTime, default=datetime.utcnow)
    updated_at = db.Column(db.DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

class ChatbotConversation(db.Model):
    __tablename__ = 'chatbot_conversations'
    
    id = db.Column(db.String(36), primary_key=True, default=lambda: str(uuid.uuid4()))
    session_id = db.Column(db.String(100), db.ForeignKey('chatbot_sessions.session_id'), nullable=False)
    query = db.Column(db.Text, nullable=False)
    response = db.Column(db.Text, nullable=False)
    confidence_score = db.Column(db.Numeric(3,2), nullable=True)
    source_documents = db.Column(db.JSON, nullable=True)
    response_time_ms = db.Column(db.Integer, nullable=True)
    created_at = db.Column(db.DateTime, default=datetime.utcnow)

class ChatbotDocument(db.Model):
    __tablename__ = 'chatbot_documents'
    
    id = db.Column(db.String(36), primary_key=True, default=lambda: str(uuid.uuid4()))
    title = db.Column(db.String(200), nullable=False)
    content = db.Column(db.Text, nullable=False)
    document_type = db.Column(db.String(50), nullable=True)
    version = db.Column(db.String(20), nullable=True)
    is_active = db.Column(db.Boolean, default=True)
    created_at = db.Column(db.DateTime, default=datetime.utcnow)
    updated_at = db.Column(db.DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

# Initialize database tables
def create_tables():
    with app.app_context():
        try:
            db.create_all()
            print("✅ Database tables created successfully!")
        except Exception as e:
            print(f"❌ Error creating tables: {e}")
            print("⚠️ Make sure PostgreSQL is running and the database 'foodflow' exists")

# Chat processing function using FoodFlow knowledge
def process_query_with_rag(query):
    """
    Process query using FoodFlow knowledge base
    """
    from src.helper import generate_response
    
    # Generate response using FoodFlow knowledge
    response, confidence = generate_response(query, "")
    sources = []
    
    return response, confidence, sources

# Routes
@app.route('/')
def index():
    # Get or create session
    if 'session_id' not in session:
        session['session_id'] = str(uuid.uuid4())
    
    # Get user_id from query parameter (optional)
    user_id = request.args.get('user_id')
    
    # Create or get chatbot session
    chatbot_session = ChatbotSession.query.filter_by(session_id=session['session_id']).first()
    if not chatbot_session:
        chatbot_session = ChatbotSession(
            user_id=user_id,
            session_id=session['session_id']
        )
        db.session.add(chatbot_session)
        db.session.commit()
    
    return render_template('chat.html', session_id=session['session_id'])

@app.route('/api/chat', methods=['POST'])
def chat():
    try:
        data = request.get_json()
        query = data.get('query', '').strip()
        session_id = data.get('session_id', session.get('session_id'))
        
        if not query:
            return jsonify({'error': 'Query is required'}), 400
        
        if not session_id:
            return jsonify({'error': 'Session ID is required'}), 400
        
        # Get or create session
        chatbot_session = ChatbotSession.query.filter_by(session_id=session_id).first()
        if not chatbot_session:
            chatbot_session = ChatbotSession(session_id=session_id)
            db.session.add(chatbot_session)
            db.session.commit()
        
        # Process query
        start_time = datetime.utcnow()
        response, confidence, sources = process_query_with_rag(query)
        response_time = (datetime.utcnow() - start_time).total_seconds() * 1000
        
        # Save conversation
        conversation = ChatbotConversation(
            session_id=session_id,
            query=query,
            response=response,
            confidence_score=confidence,
            source_documents=sources,
            response_time_ms=int(response_time)
        )
        db.session.add(conversation)
        db.session.commit()
        
        return jsonify({
            'response': response,
            'confidence': confidence,
            'session_id': session_id,
            'sources': sources,
            'response_time_ms': int(response_time)
        })
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/history/<session_id>')
def get_history(session_id):
    try:
        conversations = ChatbotConversation.query.filter_by(session_id=session_id).order_by(ChatbotConversation.created_at).all()
        
        history = []
        for conv in conversations:
            history.append({
                'id': conv.id,
                'query': conv.query,
                'response': conv.response,
                'confidence': float(conv.confidence_score) if conv.confidence_score else None,
                'timestamp': conv.created_at.isoformat(),
                'response_time_ms': conv.response_time_ms
            })
        
        return jsonify({'history': history})
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/health')
def health_check():
    try:
        # Test database connection
        from sqlalchemy import text
        db.session.execute(text('SELECT 1'))
        return jsonify({
            'status': 'healthy',
            'database': 'connected',
            'timestamp': datetime.utcnow().isoformat()
        })
    except Exception as e:
        return jsonify({
            'status': 'unhealthy',
            'database': 'disconnected',
            'error': str(e),
            'timestamp': datetime.utcnow().isoformat()
        }), 500

if __name__ == '__main__':
    print("🚀 Starting FoodFlow Chatbot...")
    print(f"📊 Database: {DATABASE_URL}")
    
    # Create tables
    create_tables()
    
    # Start Flask app
    print("🌐 Starting Flask app on http://localhost:5001")
    app.run(host='0.0.0.0', port=5001, debug=True)
