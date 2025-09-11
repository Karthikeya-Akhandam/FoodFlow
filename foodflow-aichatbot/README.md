# FoodFlow AI Chatbot

An intelligent chatbot for the FoodFlow food donation platform that helps users with questions about features, usage, and troubleshooting.

## 🚀 Quick Start

### 1. Install Dependencies
```bash
pip install -r requirements.txt
```

### 2. Setup Environment
Copy `env.local` and update with your actual values:
```bash
# Database Configuration
DATABASE_URL=postgres://postgres:12345678@localhost:5432/foodflow?sslmode=disable

# API Keys (Add your actual keys here)
PINECONE_API_KEY=your-personal-pinecone-api-key
OPENAI_API_KEY=your-personal-openai-api-key

# PDF Processing (Add your PDF path here)
PDF_DOCUMENTATION_PATH=/path/to/your/app/documentation.pdf
```

### 3. Run the Chatbot
```bash
python run_chatbot.py
```

The chatbot will be available at: http://localhost:5001

## 🏗️ Architecture

### Database Integration
- **PostgreSQL**: Connects to your existing FoodFlow database
- **New Tables**: Adds chatbot-specific tables for conversations and sessions
- **User Integration**: Links chatbot sessions to your existing users

### Chatbot Features
- **Session Management**: Tracks conversations per user/session
- **Response Logging**: Stores all queries and responses
- **Confidence Scoring**: Tracks response confidence levels
- **Performance Monitoring**: Measures response times

## 📊 Database Schema

### New Tables Added:
```sql
-- Chatbot Sessions
chatbot_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    session_id VARCHAR(100) UNIQUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)

-- Chatbot Conversations
chatbot_conversations (
    id UUID PRIMARY KEY,
    session_id VARCHAR(100) REFERENCES chatbot_sessions(session_id),
    query TEXT,
    response TEXT,
    confidence_score DECIMAL(3,2),
    source_documents JSONB,
    response_time_ms INTEGER,
    created_at TIMESTAMP
)

-- Chatbot Documents (for future PDF processing)
chatbot_documents (
    id UUID PRIMARY KEY,
    title VARCHAR(200),
    content TEXT,
    document_type VARCHAR(50),
    version VARCHAR(20),
    is_active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
```

## 🔧 API Endpoints

### Chat Endpoint
```http
POST /api/chat
Content-Type: application/json

{
    "query": "How do I create a donation offer?",
    "session_id": "optional-session-id"
}
```

### Chat History
```http
GET /api/history/{session_id}
```

### Health Check
```http
GET /api/health
```

## 🎯 Integration with Your Main App

### Frontend Integration
Add a link in your main app to redirect to the chatbot:
```html
<a href="http://localhost:5001?user_id={{ current_user.id }}" target="_blank">
    <i class="fas fa-robot"></i> Ask AI Assistant
</a>
```

### Backend Integration (Optional)
```python
import requests

def get_chatbot_response(query, user_id):
    response = requests.post('http://localhost:5001/api/chat', json={
        'query': query,
        'user_id': user_id
    })
    return response.json()
```

## 🧠 Code-Based Knowledge System

The chatbot uses your actual FoodFlow codebase to provide intelligent responses:

1. **Codebase Analysis**: Automatically extracts knowledge from your Go backend and Next.js frontend
2. **API Documentation**: Built-in knowledge of all your API endpoints and routes
3. **Database Schema**: Understands your PostgreSQL database structure
4. **Business Logic**: Knows your donation flow, matching algorithm, and credit system
5. **Real-Time Accuracy**: Always up-to-date with your actual code

## 🛠️ Development

### Project Structure
```
├── app.py                 # Main Flask application
├── run_chatbot.py         # Simple runner script
├── templates/
│   └── chat.html         # Chat interface
├── src/
│   ├── helper.py         # Utility functions
│   └── prompt.py         # AI prompt templates
├── env.local             # Environment configuration
└── requirements.txt      # Dependencies
```

### Code-Based Knowledge Features
1. **Automatic Knowledge Extraction**: Analyzes your Go backend code
2. **API Endpoint Recognition**: Understands all your routes and handlers
3. **Database Schema Understanding**: Knows your PostgreSQL table structure
4. **Business Logic Intelligence**: Understands donation flow and matching algorithm
5. **Real-Time Updates**: Knowledge stays current with your code changes

## 🔒 Security & Privacy

- **Data Privacy**: Conversations are logged for quality improvement
- **User Context**: Optional user_id linking for personalization
- **Safety Disclaimers**: Built-in disclaimers for critical operations
- **Input Sanitization**: All user inputs are sanitized

## 📈 Analytics

The chatbot tracks:
- Query types and frequency
- Response confidence scores
- Response times
- User satisfaction (future feature)
- Popular questions and topics

## 🚨 Troubleshooting

### Database Connection Issues
- Verify PostgreSQL is running
- Check database credentials in `env.local`
- Ensure database exists and is accessible

### Port Conflicts
- Default port is 5001
- Change `CHATBOT_PORT` in `env.local` if needed

### Missing Dependencies
```bash
pip install -r requirements.txt
```

## 🔮 Future Enhancements

- [ ] PDF document processing
- [ ] Vector database integration
- [ ] AI response generation
- [ ] User satisfaction tracking
- [ ] Multi-language support
- [ ] Voice interface
- [ ] Mobile app integration

## 📞 Support

For issues or questions:
1. Check the troubleshooting section
2. Review the logs for error messages
3. Verify database connectivity
4. Check environment configuration

---

**Note**: This is a prototype implementation. For production use, add proper error handling, security measures, and performance optimizations.