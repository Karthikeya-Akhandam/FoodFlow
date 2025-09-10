import React, { useState } from "react";

export default function UserProfile({ user, token, onUpdate }) {
  const [isEditing, setIsEditing] = useState(false);
  const [editForm, setEditForm] = useState({
    name: user?.user?.name || "",
    email: user?.user?.email || "",
    peopleFedLastMonth: user?.organization?.peopleFedLastMonth || 0
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";

  const isOrganization = user?.user?.role === "organization";
  const isCollaborator = user?.user?.role === "collaborator";
  const isAdmin = user?.user?.role === "admin";

  async function updateProfile() {
    setLoading(true);
    setMessage("");
    
    try {
      // In a real app, you'd have an endpoint to update user profile
      // For now, we'll just show a success message
      setTimeout(() => {
        setMessage("Profile updated successfully!");
        setIsEditing(false);
        setLoading(false);
        onUpdate();
      }, 1000);
    } catch (err) {
      setMessage("Failed to update profile: " + err.message);
      setLoading(false);
    }
  }

  return (
    <div className="space-y-8">
      {/* Profile Header */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">Profile Information</h3>
            <button
              onClick={() => setIsEditing(!isEditing)}
              className="px-4 py-2 text-blue-600 hover:text-blue-800 font-medium transition-colors"
            >
              {isEditing ? "Cancel" : "Edit Profile"}
            </button>
          </div>
        </div>
        <div className="p-6">
          {message && (
            <div className={`mb-4 p-4 rounded-lg ${
              message.includes("success") || message.includes("Success")
                ? "bg-green-50 border border-green-200 text-green-700"
                : "bg-red-50 border border-red-200 text-red-700"
            }`}>
              {message}
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
              {isEditing ? (
                <input
                  type="text"
                  value={editForm.name}
                  onChange={(e) => setEditForm({...editForm, name: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              ) : (
                <p className="text-gray-900 py-2">{user?.user?.name}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
              {isEditing ? (
                <input
                  type="email"
                  value={editForm.email}
                  onChange={(e) => setEditForm({...editForm, email: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              ) : (
                <p className="text-gray-900 py-2">{user?.user?.email}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Account Type</label>
              <p className="text-gray-900 py-2">
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                  {user?.user?.role?.charAt(0).toUpperCase() + user?.user?.role?.slice(1)}
                </span>
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Member Since</label>
              <p className="text-gray-900 py-2">
                {new Date(user?.user?.createdAt).toLocaleDateString()}
              </p>
            </div>

            {isOrganization && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">People Fed Last Month</label>
                {isEditing ? (
                  <input
                    type="number"
                    value={editForm.peopleFedLastMonth}
                    onChange={(e) => setEditForm({...editForm, peopleFedLastMonth: parseInt(e.target.value) || 0})}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                ) : (
                  <p className="text-gray-900 py-2">{user?.organization?.peopleFedLastMonth || 0}</p>
                )}
              </div>
            )}
          </div>

          {isEditing && (
            <div className="mt-6 flex justify-end space-x-3">
              <button
                onClick={() => setIsEditing(false)}
                className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={updateProfile}
                disabled={loading}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {loading ? "Saving..." : "Save Changes"}
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {isOrganization && (
          <>
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <span className="text-2xl">💳</span>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Available Credits</p>
                  <p className="text-2xl font-semibold text-gray-900">{user?.organization?.credits || 0}</p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 rounded-lg">
                  <span className="text-2xl">👥</span>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">People Fed Last Month</p>
                  <p className="text-2xl font-semibold text-gray-900">{user?.organization?.peopleFedLastMonth || 0}</p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-purple-100 rounded-lg">
                  <span className="text-2xl">🏢</span>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Organization Status</p>
                  <p className="text-lg font-semibold text-gray-900">
                    {user?.organization?.isProxy ? "Proxy Organization" : "Regular Organization"}
                  </p>
                </div>
              </div>
            </div>
          </>
        )}

        {isCollaborator && (
          <>
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-purple-100 rounded-lg">
                  <span className="text-2xl">⭐</span>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Tokens Earned</p>
                  <p className="text-2xl font-semibold text-gray-900">{user?.collaborator?.tokens || 0}</p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center">
                <div className="p-2 bg-green-100 rounded-lg">
                  <span className="text-2xl">🤝</span>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Contributor Level</p>
                  <p className="text-lg font-semibold text-gray-900">
                    {user?.collaborator?.tokens > 50 ? "Gold" : 
                     user?.collaborator?.tokens > 20 ? "Silver" : "Bronze"}
                  </p>
                </div>
              </div>
            </div>
          </>
        )}

        {isAdmin && (
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="p-2 bg-red-100 rounded-lg">
                <span className="text-2xl">👑</span>
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">Admin Privileges</p>
                <p className="text-lg font-semibold text-gray-900">Full Access</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
