#!/usr/bin/env python3
"""
FoodFlow Chatbot Runner
Simple script to start the chatbot server
"""

import os
import sys
from pathlib import Path

# Add src to path
sys.path.append(str(Path(__file__).parent / "src"))

def main():
    """Main function to run the chatbot"""
    print("🚀 Starting FoodFlow Chatbot...")
    print("=" * 50)
    
    # Check if .env file exists
    env_file = Path("env.local")
    if env_file.exists():
        print("✅ Found environment file: env.local")
    else:
        print("⚠️ No environment file found. Using default settings.")
    
    # Check database connection
    try:
        from app import app, db
        with app.app_context():
            try:
                from sqlalchemy import text
                db.session.execute(text('SELECT 1'))
                print("✅ Database connection successful")
            except Exception as db_error:
                print(f"❌ Database connection failed: {db_error}")
                print("⚠️ Continuing anyway - database will be created when needed")
    except Exception as e:
        print(f"❌ Error checking database: {e}")
        print("⚠️ Continuing anyway - database will be created when needed")
    
    # Start the Flask app
    try:
        print("🌐 Starting Flask app on http://localhost:5001")
        print("📱 Open your browser and go to: http://localhost:5001")
        print("🔄 Press Ctrl+C to stop the server")
        print("=" * 50)
        
        app.run(
            host='0.0.0.0',
            port=5001,
            debug=True
        )
    except KeyboardInterrupt:
        print("\n👋 Chatbot server stopped by user")
    except Exception as e:
        print(f"❌ Error starting server: {e}")

if __name__ == "__main__":
    main()
