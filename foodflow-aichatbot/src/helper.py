"""
Helper functions for the FoodFlow Chatbot
"""
import os
import logging
from datetime import datetime
from typing import List, Dict, Any, Optional
import json

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def setup_environment():
    """Setup environment variables and configuration"""
    try:
        from dotenv import load_dotenv
        load_dotenv()
        logger.info("✅ Environment variables loaded successfully")
    except ImportError:
        logger.warning("⚠️ python-dotenv not installed, using system environment variables")
    except Exception as e:
        logger.error(f"❌ Error loading environment: {e}")

def validate_database_connection(db):
    """Validate database connection"""
    try:
        db.session.execute('SELECT 1')
        logger.info("✅ Database connection successful")
        return True
    except Exception as e:
        logger.error(f"❌ Database connection failed: {e}")
        return False

def format_response_time(response_time_ms: int) -> str:
    """Format response time for display"""
    if response_time_ms < 1000:
        return f"{response_time_ms}ms"
    else:
        return f"{response_time_ms/1000:.1f}s"

def sanitize_input(text: str) -> str:
    """Sanitize user input"""
    if not text:
        return ""
    
    # Remove potentially harmful characters
    sanitized = text.strip()
    # Add more sanitization as needed
    
    return sanitized

def create_session_id() -> str:
    """Create a unique session ID"""
    import uuid
    return str(uuid.uuid4())

def log_conversation(session_id: str, query: str, response: str, confidence: float, response_time: int):
    """Log conversation for analytics"""
    log_data = {
        'timestamp': datetime.utcnow().isoformat(),
        'session_id': session_id,
        'query': query[:100] + '...' if len(query) > 100 else query,  # Truncate for logging
        'response_length': len(response),
        'confidence': confidence,
        'response_time_ms': response_time
    }
    
    logger.info(f"💬 Conversation logged: {json.dumps(log_data)}")

def get_confidence_color(confidence: float) -> str:
    """Get color code based on confidence score"""
    if confidence >= 0.8:
        return "#28a745"  # Green - High confidence
    elif confidence >= 0.6:
        return "#ffc107"  # Yellow - Medium confidence
    else:
        return "#dc3545"  # Red - Low confidence

def format_confidence_score(confidence: float) -> str:
    """Format confidence score for display"""
    if confidence is None:
        return "N/A"
    
    percentage = confidence * 100
    if percentage >= 80:
        return f"{percentage:.0f}% (High)"
    elif percentage >= 60:
        return f"{percentage:.0f}% (Medium)"
    else:
        return f"{percentage:.0f}% (Low)"

def extract_keywords(query: str) -> List[str]:
    """Extract keywords from user query"""
    # Simple keyword extraction
    stop_words = {'the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by', 'is', 'are', 'was', 'were', 'be', 'been', 'have', 'has', 'had', 'do', 'does', 'did', 'will', 'would', 'could', 'should', 'may', 'might', 'can', 'how', 'what', 'when', 'where', 'why', 'who'}
    
    words = query.lower().split()
    keywords = [word for word in words if word not in stop_words and len(word) > 2]
    
    return keywords

def categorize_query(query: str) -> str:
    """Categorize the type of query"""
    query_lower = query.lower()
    
    if any(word in query_lower for word in ['how', 'tutorial', 'guide', 'step']):
        return 'tutorial'
    elif any(word in query_lower for word in ['what', 'feature', 'function']):
        return 'feature_info'
    elif any(word in query_lower for word in ['error', 'problem', 'issue', 'troubleshoot']):
        return 'troubleshooting'
    elif any(word in query_lower for word in ['api', 'endpoint', 'integration']):
        return 'api_help'
    else:
        return 'general'

def prepare_rag_context(query: str, documents: List[Dict]) -> str:
    """Prepare context for RAG system"""
    if not documents:
        return ""
    
    context_parts = []
    for doc in documents:
        context_parts.append(f"Source: {doc.get('title', 'Unknown')}\n{doc.get('content', '')}")
    
    return "\n\n".join(context_parts)

def validate_pdf_file(file_path: str) -> bool:
    """Validate PDF file"""
    if not os.path.exists(file_path):
        logger.error(f"❌ PDF file not found: {file_path}")
        return False
    
    if not file_path.lower().endswith('.pdf'):
        logger.error(f"❌ File is not a PDF: {file_path}")
        return False
    
    # Check file size (max 50MB)
    file_size = os.path.getsize(file_path)
    if file_size > 50 * 1024 * 1024:  # 50MB
        logger.error(f"❌ PDF file too large: {file_size / (1024*1024):.1f}MB")
        return False
    
    logger.info(f"✅ PDF file validated: {file_path} ({file_size / (1024*1024):.1f}MB)")
    return True

def chunk_text(text: str, chunk_size: int = 1000, overlap: int = 200) -> List[str]:
    """Split text into overlapping chunks"""
    if len(text) <= chunk_size:
        return [text]
    
    chunks = []
    start = 0
    
    while start < len(text):
        end = start + chunk_size
        
        # Try to break at sentence boundary
        if end < len(text):
            # Look for sentence endings
            for i in range(end, max(start + chunk_size - 100, start), -1):
                if text[i] in '.!?':
                    end = i + 1
                    break
        
        chunk = text[start:end].strip()
        if chunk:
            chunks.append(chunk)
        
        start = end - overlap
    
    return chunks

def clean_text(text: str) -> str:
    """Clean and normalize text"""
    import re
    
    # Remove extra whitespace
    text = re.sub(r'\s+', ' ', text)
    
    # Remove special characters but keep basic punctuation
    text = re.sub(r'[^\w\s\.\,\!\?\;\:\-\(\)]', '', text)
    
    # Remove multiple periods
    text = re.sub(r'\.{2,}', '.', text)
    
    return text.strip()

# Placeholder functions for future PDF processing
def process_pdf_document(file_path: str) -> List[Dict[str, Any]]:
    """
    Process PDF document and extract text chunks
    TODO: Implement actual PDF processing
    """
    logger.info(f"📄 Processing PDF: {file_path}")
    
    # Placeholder implementation
    # In the future, this will:
    # 1. Extract text from PDF
    # 2. Clean and chunk the text
    # 3. Create embeddings
    # 4. Store in vector database
    
    return []

def search_documents(query: str, limit: int = 5) -> List[Dict[str, Any]]:
    """
    Search for relevant documents
    TODO: Implement actual vector search
    """
    logger.info(f"🔍 Searching documents for: {query}")
    
    # Placeholder implementation
    # In the future, this will:
    # 1. Create query embedding
    # 2. Search vector database
    # 3. Return relevant documents
    
    return []

def generate_response(query: str, context: str) -> tuple[str, float]:
    """
    Generate personalized response using FoodFlow knowledge base
    """
    logger.info(f"🤖 Generating personalized FoodFlow response for: {query}")
    
    # Enhanced FoodFlow knowledge base with more detailed information
    foodflow_knowledge = {
        'user_roles': {
            'ADMIN': {
                'description': 'System administrators who manage the platform and oversee operations',
                'permissions': ['User management', 'System configuration', 'Analytics access', 'Content moderation'],
                'access_level': 'Full system access'
            },
            'ORG': {
                'description': 'Organizations that receive food donations (NGOs, shelters, community centers)',
                'permissions': ['Browse offers', 'Make claims', 'Manage credits', 'View history'],
                'access_level': 'Organization-focused features'
            },
            'COLLAB': {
                'description': 'Collaborators who donate food (restaurants, hotels, caterers, individual donors)',
                'permissions': ['Create offers', 'Manage donations', 'Earn tokens', 'View analytics'],
                'access_level': 'Donation-focused features'
            }
        },
        'api_endpoints': {
            'auth': {
                'endpoints': ['POST /v1/auth/signup', 'POST /v1/auth/login', 'GET /v1/auth/me'],
                'description': 'User authentication and account management'
            },
            'offers': {
                'endpoints': ['POST /v1/offers', 'GET /v1/offers', 'GET /v1/offers/nearby'],
                'description': 'Donation offer creation and browsing'
            },
            'claims': {
                'endpoints': ['POST /v1/claims', 'GET /v1/claims', 'PUT /v1/claims/:id/status'],
                'description': 'Claim management and status updates'
            },
            'credits': {
                'endpoints': ['GET /v1/credits', 'GET /v1/credits/history', 'POST /v1/credits/spend'],
                'description': 'Credit system for organizations'
            },
            'tokens': {
                'endpoints': ['GET /v1/tokens', 'GET /v1/tokens/history', 'POST /v1/tokens/redeem'],
                'description': 'Token system for collaborators'
            }
        },
        'credit_system': {
            'allocation': 'Monthly credits are allocated to organizations based on their size and needs',
            'calculation': 'Base credits + remote bonus (if applicable) + special circumstances',
            'usage': '1 credit = 1 serving of food',
            'expiry': 'Credits expire at the end of each month',
            'bonus': 'Remote organizations get 20% bonus credits',
            'pressure': 'Organizations with low remaining credits get higher priority in matching'
        },
        'matching_algorithm': {
            'proximity': '40% - Geographic distance between collaborator and organization',
            'purpose': '25% - Alignment between food purpose and organization focus',
            'credit_pressure': '20% - Organizations with fewer remaining credits get priority',
            'reliability': '10% - Historical performance of both parties',
            'remote_bonus': '5% - Additional points for remote organization support'
        },
        'workflows': {
            'donation_flow': [
                '1. Collaborator creates donation offer with details',
                '2. System calculates priority scores for nearby organizations',
                '3. Organizations browse and claim offers',
                '4. System matches best organization to offer',
                '5. Food transfer is completed and recorded',
                '6. Credits are spent and tokens are earned'
            ],
            'registration_flow': [
                '1. User signs up with email and password',
                '2. User selects role (ORG/COLLAB)',
                '3. User completes profile with organization/collaborator details',
                '4. Account is activated and ready to use'
            ]
        }
    }
    
    query_lower = query.lower()
    keywords = extract_keywords(query)
    
    # More sophisticated response generation
    if any(word in query_lower for word in ['role', 'admin', 'org', 'collab', 'user type']):
        response = "**FoodFlow User Roles & Permissions:**\n\n"
        for role, details in foodflow_knowledge['user_roles'].items():
            response += f"**{role}** - {details['description']}\n"
            response += f"• **Access Level**: {details['access_level']}\n"
            response += f"• **Key Permissions**: {', '.join(details['permissions'])}\n\n"
        response += "Each role is designed for specific use cases in the food donation ecosystem."
        confidence = 0.95
        
    elif any(word in query_lower for word in ['credit', 'credits']):
        if any(word in query_lower for word in ['how', 'calculate', 'awarded', 'given']):
            response = "**How Credits Work in FoodFlow:**\n\n"
            response += f"**Allocation**: {foodflow_knowledge['credit_system']['allocation']}\n"
            response += f"**Calculation**: {foodflow_knowledge['credit_system']['calculation']}\n"
            response += f"**Usage**: {foodflow_knowledge['credit_system']['usage']}\n"
            response += f"**Expiry**: {foodflow_knowledge['credit_system']['expiry']}\n"
            response += f"**Remote Bonus**: {foodflow_knowledge['credit_system']['bonus']}\n\n"
            response += "**Credit Pressure**: Organizations with fewer credits get higher priority in the matching algorithm."
            confidence = 0.95
        elif any(word in query_lower for word in ['run out', 'finish', 'expire', 'end']):
            response = "**What Happens When Credits Run Out:**\n\n"
            response += "• **Priority Boost**: Organizations with low credits get higher priority in matching\n"
            response += "• **Monthly Reset**: New credits are allocated at the start of each month\n"
            response += "• **Remote Support**: Remote organizations get bonus credits to ensure fair access\n"
            response += "• **Emergency Access**: Special provisions may be available for urgent needs\n\n"
            response += "The system is designed to ensure organizations can always access food when needed."
            confidence = 0.9
        else:
            response = "**FoodFlow Credit System:**\n\n"
            response += f"**Purpose**: {foodflow_knowledge['credit_system']['allocation']}\n"
            response += f"**Value**: {foodflow_knowledge['credit_system']['usage']}\n"
            response += f"**Timeline**: {foodflow_knowledge['credit_system']['expiry']}\n"
            response += f"**Special**: {foodflow_knowledge['credit_system']['bonus']}\n\n"
            response += "Credits ensure fair distribution and prevent abuse of the donation system."
            confidence = 0.9
            
    elif any(word in query_lower for word in ['token', 'tokens']):
        response = "**FoodFlow Token System:**\n\n"
        response += "**For Collaborators (Food Donors):**\n"
        response += "• **Earning**: 1 token per successful donation\n"
        response += "• **Redemption**: Can be exchanged for credits or rewards\n"
        response += "• **Tracking**: Monthly history and analytics available\n"
        response += "• **Incentive**: Encourages consistent food donations\n\n"
        response += "**Token Benefits:**\n"
        response += "• Recognition for contribution to food security\n"
        response += "• Potential rewards and partnerships\n"
        response += "• Priority in special programs\n"
        response += "• Community impact tracking"
        confidence = 0.9
        
    elif any(word in query_lower for word in ['matching', 'algorithm', 'priority', 'score']):
        response = "**FoodFlow Matching Algorithm:**\n\n"
        response += "The system uses a sophisticated scoring system to match offers with organizations:\n\n"
        for factor, weight in foodflow_knowledge['matching_algorithm'].items():
            response += f"• **{factor.replace('_', ' ').title()}**: {weight}\n"
        response += "\n**How It Works:**\n"
        response += "1. System calculates scores for all eligible organizations\n"
        response += "2. Organizations are ranked by total score\n"
        response += "3. Highest-scoring organization gets the offer\n"
        response += "4. If declined, next organization in line gets priority"
        confidence = 0.95
        
    elif any(word in query_lower for word in ['api', 'endpoint', 'route', 'integration']):
        response = "**FoodFlow API Endpoints:**\n\n"
        for category, details in foodflow_knowledge['api_endpoints'].items():
            response += f"**{category.title()}** - {details['description']}\n"
            for endpoint in details['endpoints']:
                response += f"• {endpoint}\n"
            response += "\n"
        response += "**Authentication**: All endpoints require JWT Bearer token except signup/login\n"
        response += "**Base URL**: http://localhost:8080\n"
        response += "**Version**: v1"
        confidence = 0.9
        
    elif any(word in query_lower for word in ['donation', 'offer', 'create', 'post']):
        response = "**Creating Donation Offers:**\n\n"
        response += "**Step-by-Step Process:**\n"
        response += "1. **Authentication**: Login as COLLAB role\n"
        response += "2. **API Call**: POST /v1/offers\n"
        response += "3. **Required Data**:\n"
        response += "   • title (5-100 characters)\n"
        response += "   • ready_from (timestamp)\n"
        response += "   • expires_at (timestamp)\n"
        response += "   • estimated_servings (1-10000)\n"
        response += "4. **Optional Data**:\n"
        response += "   • description\n"
        response += "   • purpose (CHILDREN, ELDERLY, WOMEN, GENERAL, EMERGENCY)\n"
        response += "   • location (pincode, city, state)\n\n"
        response += "**Validation Rules**:\n"
        response += "• expires_at must be after ready_from\n"
        response += "• servings must be between 1-10000\n"
        response += "• title must be meaningful and descriptive"
        confidence = 0.95
        
    elif any(word in query_lower for word in ['claim', 'organization', 'request']):
        response = "**Making Claims on Offers:**\n\n"
        response += "**For Organizations (ORG role):**\n"
        response += "1. **Browse Offers**: GET /v1/offers/nearby\n"
        response += "2. **Review Details**: Check serving size, location, purpose\n"
        response += "3. **Make Claim**: POST /v1/claims\n"
        response += "4. **Required Data**:\n"
        response += "   • offer_id\n"
        response += "   • requested_servings\n"
        response += "5. **Wait for Matching**: System processes all claims\n"
        response += "6. **Get Notification**: If selected, confirm pickup\n\n"
        response += "**Matching Process**:\n"
        response += "• All claims are scored using the matching algorithm\n"
        response += "• Highest-scoring organization wins the offer\n"
        response += "• Others are notified and can try other offers"
        confidence = 0.95
        
    elif any(word in query_lower for word in ['workflow', 'process', 'flow', 'how it works']):
        response = "**FoodFlow Complete Workflow:**\n\n"
        response += "**Donation Process:**\n"
        for step in foodflow_knowledge['workflows']['donation_flow']:
            response += f"{step}\n"
        response += "\n**Registration Process:**\n"
        for step in foodflow_knowledge['workflows']['registration_flow']:
            response += f"{step}\n"
        response += "\n**Key Features:**\n"
        response += "• Real-time matching and notifications\n"
        response += "• Credit-based fair distribution\n"
        response += "• Geographic proximity optimization\n"
        response += "• Remote organization support"
        confidence = 0.9
        
    elif any(word in query_lower for word in ['database', 'table', 'schema', 'data']):
        response = "**FoodFlow Database Schema:**\n\n"
        response += "**Core Tables:**\n"
        response += "• **users** - Authentication and role management\n"
        response += "• **profiles** - User profile information\n"
        response += "• **organizations** - NGO/shelter details\n"
        response += "• **collaborators** - Restaurant/hotel details\n"
        response += "• **donation_offers** - Food donation listings\n"
        response += "• **donation_claims** - Organization requests\n"
        response += "• **redemptions** - Completed transactions\n"
        response += "• **credits** - Monthly credit allocations\n"
        response += "• **tokens** - Collaborator rewards\n"
        response += "• **events** - Complete audit trail\n\n"
        response += "**Technology**: PostgreSQL with UUID primary keys and JSONB for flexible data"
        confidence = 0.9
        
    else:
        # Personalized response based on query content
        response = f"Based on your question about **'{query}'**, here's what I can help you with:\n\n"
        
        if any(word in keywords for word in ['help', 'support', 'assistance']):
            response += "**I can help you with:**\n"
            response += "• Understanding user roles and permissions\n"
            response += "• Explaining the credit and token systems\n"
            response += "• API endpoint documentation\n"
            response += "• Database schema information\n"
            response += "• Workflow and process guidance\n"
            response += "• Matching algorithm details\n\n"
            response += "**Just ask me specific questions like:**\n"
            response += "• 'How do credits work?'\n"
            response += "• 'What are the user roles?'\n"
            response += "• 'How does the matching algorithm work?'"
            confidence = 0.8
        else:
            response += "**FoodFlow Overview:**\n"
            response += "FoodFlow is a smart food donation platform that connects restaurants and hotels with NGOs and shelters to reduce food waste and help those in need.\n\n"
            response += "**Key Features:**\n"
            response += "• **Smart Matching**: AI-powered algorithm matches offers with organizations\n"
            response += "• **Credit System**: Fair distribution through monthly credit allocations\n"
            response += "• **Token Rewards**: Incentivizes consistent food donations\n"
            response += "• **Geographic Optimization**: Prioritizes local matches\n"
            response += "• **Remote Support**: Special provisions for remote organizations\n\n"
            response += "**Ask me about:**\n"
            response += "• User roles (ADMIN, ORG, COLLAB)\n"
            response += "• Credit and token systems\n"
            response += "• API endpoints and integration\n"
            response += "• Database structure\n"
            response += "• Matching algorithm\n"
            response += "• Workflow processes"
            confidence = 0.7
    
    return response, confidence
