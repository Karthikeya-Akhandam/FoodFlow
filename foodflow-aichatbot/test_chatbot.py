#!/usr/bin/env python3
"""
Test script for FoodFlow Chatbot
"""
import sys
from pathlib import Path

# Add src to path
sys.path.append(str(Path(__file__).parent / "src"))

def test_foodflow_knowledge():
    """Test the FoodFlow knowledge system"""
    print("🧪 Testing FoodFlow Chatbot Knowledge System")
    print("=" * 50)
    
    from src.helper import generate_response
    
    # Test queries
    test_queries = [
        "What are the user roles in FoodFlow?",
        "How do I create a donation offer?",
        "What API endpoints are available?",
        "How does the credit system work?",
        "What database tables exist?",
        "How do organizations make claims?",
        "What is the matching algorithm?",
        "Tell me about FoodFlow"
    ]
    
    for i, query in enumerate(test_queries, 1):
        print(f"\n{i}. Query: {query}")
        print("-" * 40)
        
        try:
            response, confidence = generate_response(query, "")
            print(f"Response (Confidence: {confidence:.2f}):")
            print(response)
        except Exception as e:
            print(f"❌ Error: {e}")
    
    print("\n" + "=" * 50)
    print("✅ FoodFlow knowledge system test completed!")

if __name__ == "__main__":
    test_foodflow_knowledge()
