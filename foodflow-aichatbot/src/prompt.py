"""
Prompt templates for the FoodFlow Chatbot
"""
from typing import Dict, List, Any

# System prompts for different types of queries
SYSTEM_PROMPTS = {
    'general': """
You are a helpful AI assistant for the FoodFlow application. FoodFlow is a food donation platform that connects organizations with food donors (restaurants, hotels, etc.) to reduce food waste and help those in need.

Your role is to:
1. Answer questions about FoodFlow features and functionality
2. Provide guidance on how to use the application
3. Help troubleshoot common issues
4. Explain the donation process and user roles

Always be helpful, accurate, and professional. If you're unsure about something, say so and suggest contacting support.
""",

    'tutorial': """
You are a tutorial assistant for FoodFlow. Help users understand how to perform specific tasks in the application.

Provide step-by-step instructions that are:
1. Clear and easy to follow
2. Specific to the FoodFlow interface
3. Include any important notes or tips
4. Mention relevant user roles (ADMIN, ORG, COLLAB)

If the task requires specific permissions or user roles, make sure to mention that.
""",

    'troubleshooting': """
You are a technical support assistant for FoodFlow. Help users resolve issues they're experiencing.

Your approach should be:
1. Ask clarifying questions if needed
2. Provide step-by-step solutions
3. Suggest common fixes first
4. Escalate to support if the issue is complex
5. Always be patient and understanding

Remember that users may be frustrated, so be empathetic and helpful.
""",

    'api_help': """
You are an API documentation assistant for FoodFlow. Help developers understand how to integrate with the FoodFlow API.

Provide information about:
1. API endpoints and their purposes
2. Request/response formats
3. Authentication requirements
4. Rate limiting and best practices
5. Example code snippets when helpful

Always mention the API version and any deprecation notices.
"""
}

# User role descriptions
USER_ROLES = {
    'ADMIN': 'System administrators who manage the platform',
    'ORG': 'Organizations that receive food donations (NGOs, shelters, etc.)',
    'COLLAB': 'Collaborators who donate food (restaurants, hotels, caterers, etc.)'
}

# Common features and their descriptions
FEATURES = {
    'donation_offers': 'Create and manage food donation offers',
    'donation_claims': 'Claim available food donations',
    'user_management': 'Manage user accounts and profiles',
    'credits_tokens': 'Credit and token system for tracking donations',
    'analytics': 'View donation statistics and reports',
    'notifications': 'Receive updates about donations and claims'
}

# Response templates
RESPONSE_TEMPLATES = {
    'greeting': "Hello! I'm the FoodFlow AI assistant. How can I help you today?",
    
    'feature_explanation': """
Based on the documentation, here's what I found about {feature}:

{explanation}

{additional_info}
""",

    'tutorial_response': """
Here's how to {action} in FoodFlow:

{steps}

{notes}
""",

    'troubleshooting_response': """
I understand you're having trouble with {issue}. Here are some steps to try:

{solutions}

If these don't work, please contact our support team for further assistance.
""",

    'no_information': """
I don't have specific information about that in the current documentation. 

However, I can help you with:
- General FoodFlow features
- User roles and permissions
- Basic troubleshooting
- Contacting support

Would you like me to help with any of these instead?
""",

    'confidence_low': """
I found some information, but I'm not completely confident in my response. Here's what I know:

{response}

For the most accurate information, I recommend:
1. Checking the official documentation
2. Contacting FoodFlow support
3. Asking a FoodFlow administrator
"""
}

# Safety and disclaimer messages
SAFETY_MESSAGES = {
    'disclaimer': """
⚠️ **Important Disclaimer:**
This AI assistant provides information based on the available documentation. For:
- Critical system issues
- Data security concerns  
- Account problems
- Payment or transaction issues

Please contact the FoodFlow support team directly.
""",

    'emergency': """
🚨 **Emergency Support:**
If you're experiencing a critical issue that affects:
- User safety
- Data security
- System availability
- Financial transactions

Please contact support immediately at [support-email] or call [support-phone].
""",

    'data_privacy': """
🔒 **Data Privacy:**
This conversation is logged for quality improvement purposes. No personal or sensitive data is stored. Your privacy is important to us.
"""
}

def get_system_prompt(query_type: str = 'general') -> str:
    """Get system prompt based on query type"""
    return SYSTEM_PROMPTS.get(query_type, SYSTEM_PROMPTS['general'])

def get_response_template(template_type: str) -> str:
    """Get response template"""
    return RESPONSE_TEMPLATES.get(template_type, RESPONSE_TEMPLATES['no_information'])

def format_feature_response(feature: str, explanation: str, additional_info: str = "") -> str:
    """Format feature explanation response"""
    return RESPONSE_TEMPLATES['feature_explanation'].format(
        feature=feature,
        explanation=explanation,
        additional_info=additional_info
    )

def format_tutorial_response(action: str, steps: str, notes: str = "") -> str:
    """Format tutorial response"""
    return RESPONSE_TEMPLATES['tutorial_response'].format(
        action=action,
        steps=steps,
        notes=notes
    )

def format_troubleshooting_response(issue: str, solutions: str) -> str:
    """Format troubleshooting response"""
    return RESPONSE_TEMPLATES['troubleshooting_response'].format(
        issue=issue,
        solutions=solutions
    )

def get_user_role_info(role: str) -> str:
    """Get information about a user role"""
    return USER_ROLES.get(role.upper(), f"Unknown role: {role}")

def get_feature_info(feature: str) -> str:
    """Get information about a feature"""
    return FEATURES.get(feature, f"Unknown feature: {feature}")

def add_safety_disclaimer(response: str, confidence: float) -> str:
    """Add appropriate safety disclaimers based on confidence"""
    if confidence < 0.6:
        response += "\n\n" + SAFETY_MESSAGES['disclaimer']
    
    return response

def create_context_prompt(query: str, context_documents: List[Dict[str, Any]]) -> str:
    """Create a prompt with context from retrieved documents"""
    if not context_documents:
        return f"Query: {query}\n\nNo relevant documentation found."
    
    context_parts = []
    for i, doc in enumerate(context_documents, 1):
        context_parts.append(f"Document {i}: {doc.get('title', 'Unknown')}\n{doc.get('content', '')}")
    
    context = "\n\n".join(context_parts)
    
    return f"""
Context from FoodFlow documentation:

{context}

Query: {query}

Please provide a helpful response based on the above context. If the context doesn't contain enough information to answer the query, say so and suggest alternative ways to get help.
"""

def create_rag_prompt(query: str, context: str, user_role: str = None) -> str:
    """Create a comprehensive RAG prompt"""
    base_prompt = f"""
You are a helpful AI assistant for FoodFlow, a food donation platform.

Context from documentation:
{context}

User Query: {query}
"""
    
    if user_role:
        base_prompt += f"\nUser Role: {user_role}\n"
    
    base_prompt += """
Please provide a helpful, accurate response based on the context. If the context doesn't contain enough information, be honest about it and suggest how the user can get more help.

Guidelines:
1. Be specific and actionable
2. Mention relevant user roles if applicable
3. Include safety disclaimers for critical operations
4. Keep responses concise but complete
5. Use a friendly, professional tone
"""
    
    return base_prompt

# Common questions and their expected response types
COMMON_QUESTIONS = {
    'how_to_create_offer': 'tutorial',
    'how_to_claim_donation': 'tutorial',
    'user_roles': 'feature_info',
    'credits_system': 'feature_info',
    'api_integration': 'api_help',
    'login_issues': 'troubleshooting',
    'password_reset': 'tutorial',
    'account_suspension': 'troubleshooting'
}

def get_expected_response_type(query: str) -> str:
    """Determine expected response type based on query"""
    query_lower = query.lower()
    
    for keyword, response_type in COMMON_QUESTIONS.items():
        if keyword.replace('_', ' ') in query_lower:
            return response_type
    
    # Default classification based on keywords
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
