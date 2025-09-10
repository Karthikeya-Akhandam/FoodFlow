import React, { useState } from "react";

export default function Login({ onLogin }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isSignup, setIsSignup] = useState(false);
  const [name, setName] = useState("");
  const [role, setRole] = useState("collaborator");
  const [peopleFed, setPeopleFed] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

  async function submit(e) {
    e.preventDefault();
    setLoading(true);
    setError("");
    
    try {
      const endpoint = isSignup ? "/auth/signup" : "/auth/login";
      const body = isSignup 
        ? { name, email, password, role, people_fed_last_month: peopleFed }
        : { email, password };
        
      const res = await fetch(`${API_BASE}${endpoint}`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(body),
      });
      const j = await res.json();
      if (res.ok) {
        onLogin(j.token);
      } else {
        setError(j.error || "Authentication failed");
      }
    } catch (err) {
      setError("Network error. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-400 via-blue-500 to-purple-600 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md overflow-hidden">
        <div className="bg-gradient-to-r from-green-500 to-blue-600 p-6 text-white text-center">
          <h1 className="text-3xl font-bold mb-2">🍽️ FoodShare</h1>
          <p className="text-green-100">Connecting communities through food</p>
        </div>
        
        <form onSubmit={submit} className="p-6 space-y-4">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
              {error}
            </div>
          )}
          
          <h2 className="text-2xl font-semibold text-gray-800 text-center mb-6">
            {isSignup ? "Create Account" : "Welcome Back"}
          </h2>
          
          {isSignup && (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
                <input 
                  value={name} 
                  onChange={e=>setName(e.target.value)} 
                  placeholder="Enter your full name" 
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all" 
                  required
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Account Type</label>
                <select 
                  value={role} 
                  onChange={e=>setRole(e.target.value)} 
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                >
                  <option value="collaborator">🤝 Collaborator (Donate Food)</option>
                  <option value="organization">🏢 Organization (Receive Food)</option>
                </select>
              </div>
              
              {role === "organization" && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">People Fed Last Month</label>
                  <input 
                    type="number" 
                    value={peopleFed} 
                    onChange={e=>setPeopleFed(parseInt(e.target.value) || 0)} 
                    placeholder="Enter number of people fed" 
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                    min="0"
                  />
                </div>
              )}
            </>
          )}
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
            <input 
              value={email} 
              onChange={e=>setEmail(e.target.value)} 
              placeholder="Enter your email" 
              type="email"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all" 
              required
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input 
              value={password} 
              onChange={e=>setPassword(e.target.value)} 
              type="password" 
              placeholder="Enter your password" 
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all" 
              required
            />
          </div>
          
          <button 
            type="submit" 
            disabled={loading}
            className="w-full bg-gradient-to-r from-green-500 to-blue-600 text-white py-3 px-4 rounded-lg font-semibold hover:from-green-600 hover:to-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Please wait..." : (isSignup ? "Create Account" : "Sign In")}
          </button>
          
          <button 
            type="button" 
            onClick={() => setIsSignup(!isSignup)} 
            className="w-full text-blue-600 hover:text-blue-800 font-medium py-2 transition-colors"
          >
            {isSignup ? "Already have an account? Sign in" : "Need an account? Sign up"}
          </button>
        </form>
      </div>
    </div>
  );
}
