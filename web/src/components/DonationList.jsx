import React, { useState } from "react";

export default function DonationList({ donations, token, onClaimed }) {
  const [claimingId, setClaimingId] = useState(null);
  const [message, setMessage] = useState("");
  const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";
  
  async function claim(id) {
    setClaimingId(id);
    setMessage("");
    
    try {
      const res = await fetch(`${API_BASE}/donations/${id}/claim`, {
        method: "POST",
        headers: { "Authorization": token, "Content-Type": "application/json" }
      });
      const j = await res.json();
      if (res.ok) {
        setMessage("Donation claimed successfully!");
        onClaimed();
      } else {
        setMessage(j.error || "Failed to claim donation");
      }
    } catch (err) {
      setMessage("Network error. Please try again.");
    } finally {
      setClaimingId(null);
    }
  }

  function formatDate(dateString) {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString();
  }

  function formatDateTime(dateString) {
    if (!dateString) return null;
    return new Date(dateString).toLocaleString();
  }

  if (donations.length === 0) {
    return (
      <div className="text-center py-8">
        <div className="text-6xl mb-4">🍽️</div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">No Donations Available</h3>
        <p className="text-gray-600">Check back later for new food donations!</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {message && (
        <div className={`p-3 rounded-lg text-sm ${
          message.includes("success") || message.includes("Success")
            ? "bg-green-50 border border-green-200 text-green-700"
            : "bg-red-50 border border-red-200 text-red-700"
        }`}>
          {message}
        </div>
      )}

      <div className="space-y-4">
        {donations.map(d => (
          <div key={d.ID} className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="flex items-center mb-2">
                  <h4 className="text-lg font-semibold text-gray-900">{d.Title}</h4>
                  <span className="ml-2 px-2 py-1 bg-green-100 text-green-800 text-xs font-medium rounded-full">
                    Available
                  </span>
                </div>
                
                {d.Description && (
                  <p className="text-gray-600 mb-2">{d.Description}</p>
                )}
                
                <div className="flex flex-wrap gap-4 text-sm text-gray-500">
                  <div className="flex items-center">
                    <span className="mr-1">📦</span>
                    {d.Quantity} items
                  </div>
                  
                  {d.Purpose && (
                    <div className="flex items-center">
                      <span className="mr-1">📝</span>
                      {d.Purpose}
                    </div>
                  )}
                  
                  {d.ExpiryAt && (
                    <div className="flex items-center">
                      <span className="mr-1">⏰</span>
                      Expires: {formatDateTime(d.ExpiryAt)}
                    </div>
                  )}
                  
                  <div className="flex items-center">
                    <span className="mr-1">👤</span>
                    By: {d.Collaborator?.User?.name || "Anonymous"}
                  </div>
                  
                  <div className="flex items-center">
                    <span className="mr-1">📅</span>
                    Posted: {formatDate(d.CreatedAt)}
                  </div>
                </div>
              </div>
              
              <div className="ml-4 flex flex-col items-end">
                <button 
                  onClick={() => claim(d.ID)}
                  disabled={claimingId === d.ID}
                  className="px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg font-medium hover:from-blue-600 hover:to-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {claimingId === d.ID ? "Claiming..." : "🎯 Claim Donation"}
                </button>
                <p className="text-xs text-gray-500 mt-1">Cost: 1 credit</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
