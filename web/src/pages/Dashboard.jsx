import React, { useEffect, useState } from "react";
import DonationForm from "../components/DonationForm";
import DonationList from "../components/DonationList";
import AdminPanel from "../components/AdminPanel";
import UserProfile from "../components/UserProfile";

export default function Dashboard({ token, onLogout }) {
  const [donations, setDonations] = useState([]);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("dashboard");
  const API_BASE = import.meta.env.VITE_API_URL || "http://localhost:8080";
  
  const api = (path, opts={}) => fetch(`${API_BASE}${path}`, {
    ...opts,
    headers: {...(opts.headers||{}), "Authorization": token, "Content-Type":"application/json"}
  });

  async function loadUser() {
    try {
      const res = await api("/me");
      const data = await res.json();
      if (res.ok) {
        setUser(data);
      }
    } catch (err) {
      console.error("Failed to load user:", err);
    }
  }

  async function loadDonations() {
    try {
      const res = await fetch(`${API_BASE}/donations`);
      const data = await res.json();
      setDonations(data);
    } catch (err) {
      console.error("Failed to load donations:", err);
    }
  }

  async function load() {
    setLoading(true);
    await Promise.all([loadUser(), loadDonations()]);
    setLoading(false);
  }

  useEffect(() => { load(); }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading your dashboard...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">Failed to load user data</p>
          <button onClick={onLogout} className="px-4 py-2 bg-red-600 text-white rounded-lg">
            Logout
          </button>
        </div>
      </div>
    );
  }

  const isAdmin = user.user?.role === "admin";
  const isCollaborator = user.user?.role === "collaborator";
  const isOrganization = user.user?.role === "organization";

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-gray-900">🍽️ FoodShare</h1>
              <span className="ml-3 px-3 py-1 bg-blue-100 text-blue-800 text-sm font-medium rounded-full">
                {user.user?.role?.charAt(0).toUpperCase() + user.user?.role?.slice(1)}
              </span>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">Welcome, {user.user?.name}</span>
              <button 
                onClick={onLogout}
                className="px-4 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-8">
            <button
              onClick={() => setActiveTab("dashboard")}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === "dashboard"
                  ? "border-blue-500 text-blue-600"
                  : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
              }`}
            >
              Dashboard
            </button>
            <button
              onClick={() => setActiveTab("profile")}
              className={`py-4 px-1 border-b-2 font-medium text-sm ${
                activeTab === "profile"
                  ? "border-blue-500 text-blue-600"
                  : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
              }`}
            >
              Profile
            </button>
            {isAdmin && (
              <button
                onClick={() => setActiveTab("admin")}
                className={`py-4 px-1 border-b-2 font-medium text-sm ${
                  activeTab === "admin"
                    ? "border-blue-500 text-blue-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                }`}
              >
                Admin Panel
              </button>
            )}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === "dashboard" && (
          <div className="space-y-8">
            {/* Stats Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="bg-white rounded-lg shadow p-6">
                <div className="flex items-center">
                  <div className="p-2 bg-green-100 rounded-lg">
                    <span className="text-2xl">🍽️</span>
                  </div>
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-600">Available Donations</p>
                    <p className="text-2xl font-semibold text-gray-900">{donations.length}</p>
                  </div>
                </div>
              </div>
              
              {isOrganization && (
                <div className="bg-white rounded-lg shadow p-6">
                  <div className="flex items-center">
                    <div className="p-2 bg-blue-100 rounded-lg">
                      <span className="text-2xl">💳</span>
                    </div>
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-600">Credits Available</p>
                      <p className="text-2xl font-semibold text-gray-900">{user.organization?.credits || 0}</p>
                    </div>
                  </div>
                </div>
              )}
              
              {isCollaborator && (
                <div className="bg-white rounded-lg shadow p-6">
                  <div className="flex items-center">
                    <div className="p-2 bg-purple-100 rounded-lg">
                      <span className="text-2xl">⭐</span>
                    </div>
                    <div className="ml-4">
                      <p className="text-sm font-medium text-gray-600">Tokens Earned</p>
                      <p className="text-2xl font-semibold text-gray-900">{user.collaborator?.tokens || 0}</p>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Main Content Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {isCollaborator && (
                <div className="bg-white rounded-lg shadow">
                  <div className="px-6 py-4 border-b border-gray-200">
                    <h3 className="text-lg font-medium text-gray-900">Create New Donation</h3>
                  </div>
                  <div className="p-6">
                    <DonationForm token={token} onCreated={load} />
                  </div>
                </div>
              )}
              
              <div className="bg-white rounded-lg shadow">
                <div className="px-6 py-4 border-b border-gray-200">
                  <h3 className="text-lg font-medium text-gray-900">Available Donations</h3>
                </div>
                <div className="p-6">
                  <DonationList donations={donations} token={token} onClaimed={load} />
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === "profile" && (
          <UserProfile user={user} token={token} onUpdate={load} />
        )}

        {activeTab === "admin" && isAdmin && (
          <AdminPanel token={token} onUpdate={load} />
        )}
      </main>
    </div>
  );
}
