"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { 
  Plus, 
  Search, 
  Filter, 
  MapPin, 
  Clock, 
  Users, 
  MoreVertical,
  Edit,
  Trash2,
  Eye,
  Gift,
  Heart,
  Target,
  Calendar,
  Star,
  Zap,
  Activity,
  CheckCircle,
  AlertCircle,
  Sparkles,
  TrendingUp,
  ArrowRight
} from "lucide-react";
import { useAuth } from "../../../lib/auth";
import { useOffers } from "../../../lib/hooks";

export default function OffersPage() {
  const { user } = useAuth();
  const { data: apiOffers, isLoading: offersLoading, error: offersError } = useOffers();
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [purposeFilter, setPurposeFilter] = useState("all");
  const [sortBy, setSortBy] = useState("newest");

  // Use real API data instead of mock data
  const offers = apiOffers || [];

  // Loading state
  if (offersLoading) {
    return (
      <div className="space-y-8">
        {/* Header Skeleton */}
        <div className="flex justify-between items-center">
          <div className="space-y-2">
            <div className="h-8 bg-slate-200 rounded w-48"></div>
            <div className="h-4 bg-slate-200 rounded w-32"></div>
          </div>
          <div className="h-10 bg-slate-200 rounded w-32"></div>
        </div>

        {/* Filters Skeleton */}
        <div className="flex flex-wrap gap-4">
          <div className="h-10 bg-slate-200 rounded w-64"></div>
          <div className="h-10 bg-slate-200 rounded w-32"></div>
          <div className="h-10 bg-slate-200 rounded w-32"></div>
          <div className="h-10 bg-slate-200 rounded w-32"></div>
        </div>

        {/* Offers Grid Skeleton */}
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="h-4 bg-slate-200 rounded w-3/4 mb-4"></div>
              <div className="h-3 bg-slate-200 rounded w-full mb-2"></div>
              <div className="h-3 bg-slate-200 rounded w-2/3 mb-4"></div>
              <div className="flex justify-between items-center">
                <div className="h-6 bg-slate-200 rounded w-20"></div>
                <div className="h-6 bg-slate-200 rounded w-16"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  // Error state
  if (offersError) {
    return (
      <div className="space-y-8">
        <div className="text-center py-12">
          <AlertCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-slate-900 mb-2">Unable to Load Offers</h2>
          <p className="text-slate-600 mb-6">
            {offersError.message || "There was an error loading the offers. Please try again later."}
          </p>
          <button 
            onClick={() => window.location.reload()} 
            className="btn-primary"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  const getStatusColor = (status) => {
    switch (status) {
      case "AVAILABLE": return "from-emerald-400 to-emerald-600 text-white";
      case "PENDING_CONFIRMATION": return "from-amber-400 to-amber-600 text-white";
      case "CLAIMED": return "from-blue-400 to-blue-600 text-white";
      case "EXPIRED": return "from-slate-400 to-slate-600 text-white";
      case "CANCELLED": return "from-red-400 to-red-600 text-white";
      default: return "from-slate-400 to-slate-600 text-white";
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case "AVAILABLE": return CheckCircle;
      case "PENDING_CONFIRMATION": return Clock;
      case "CLAIMED": return Users;
      case "EXPIRED": return AlertCircle;
      case "CANCELLED": return AlertCircle;
      default: return AlertCircle;
    }
  };

  const getPurposeColor = (purpose) => {
    switch (purpose) {
      case "CHILDREN": return "from-pink-400 to-pink-600 text-white";
      case "ELDERLY": return "from-purple-400 to-purple-600 text-white";
      case "WOMEN": return "from-rose-400 to-rose-600 text-white";
      case "GENERAL": return "from-blue-400 to-blue-600 text-white";
      case "EMERGENCY": return "from-red-400 to-red-600 text-white";
      default: return "from-slate-400 to-slate-600 text-white";
    }
  };

  const getPurposeIcon = (purpose) => {
    switch (purpose) {
      case "CHILDREN": return Heart;
      case "ELDERLY": return Users;
      case "WOMEN": return Heart;
      case "GENERAL": return Target;
      case "EMERGENCY": return AlertCircle;
      default: return Target;
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case "URGENT": return "from-red-400 to-red-600";
      case "HIGH": return "from-orange-400 to-orange-600";
      case "MEDIUM": return "from-yellow-400 to-yellow-600";
      case "LOW": return "from-green-400 to-green-600";
      default: return "from-slate-400 to-slate-600";
    }
  };

  const filteredOffers = offers.filter(offer => {
    const matchesSearch = (offer.title || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
                         (offer.description || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
                         (offer.organization_name || offer.collaborator_name || '').toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === "all" || offer.status === statusFilter;
    const matchesPurpose = purposeFilter === "all" || offer.purpose === purposeFilter;
    return matchesSearch && matchesStatus && matchesPurpose;
  });

  const sortedOffers = [...filteredOffers].sort((a, b) => {
    switch (sortBy) {
      case "newest": return new Date(b.created_at || b.createdAt) - new Date(a.created_at || a.createdAt);
      case "oldest": return new Date(a.created_at || a.createdAt) - new Date(b.created_at || b.createdAt);
      case "servings": return (b.servings || 0) - (a.servings || 0);
      case "distance": return parseFloat(a.distance || 0) - parseFloat(b.distance || 0);
      case "rating": return (b.rating || 0) - (a.rating || 0);
      default: return new Date(b.createdAt) - new Date(a.createdAt);
    }
  });

  const canCreateOffer = user?.role === "COLLAB" || user?.role === "ADMIN";
  const canClaimOffer = user?.role === "ORG" || user?.role === "ADMIN";

  return (
    <div className="space-y-6 lg:space-y-8">
      {/* Header Section */}
      <div className="modern-card p-6 lg:p-8 bg-gradient-to-br from-emerald-50 to-cyan-50 border-emerald-200 animate-fade-in">
        <div className="flex flex-col gap-6">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
            <div className="flex items-center space-x-3 lg:space-x-4">
              <div className="w-12 h-12 lg:w-16 lg:h-16 bg-gradient-primary rounded-xl lg:rounded-2xl flex items-center justify-center flex-shrink-0">
                <Gift className="w-6 h-6 lg:w-8 lg:h-8 text-white" />
              </div>
              <div className="min-w-0">
                <h1 className="text-2xl lg:text-4xl font-bold text-slate-900 leading-tight">
                  {user?.role === "COLLAB" ? "My Food Offers" : "Available Food Donations"}
                </h1>
                <p className="text-emerald-600 font-medium text-sm lg:text-lg">
                  {user?.role === "COLLAB" 
                    ? "Manage and track your food donation offers"
                    : "Discover and claim food donations in your area"
                  }
                </p>
              </div>
            </div>
            {canCreateOffer && (
              <div className="flex-shrink-0">
                <Link
                  href="/dashboard/offers/new"
                  className="btn-modern px-6 lg:px-8 py-3 lg:py-4 text-base lg:text-lg inline-flex items-center group w-full lg:w-auto justify-center"
                >
                  <Plus className="h-5 w-5 lg:h-6 lg:w-6 mr-2 lg:mr-3 group-hover:scale-110 transition-transform" />
                  Create New Offer
                </Link>
              </div>
            )}
          </div>
          
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-slate-900">{offers.length}</div>
              <div className="text-xs lg:text-sm text-slate-600">Total Offers</div>
            </div>
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-emerald-600">{offers.filter(o => o.status === "AVAILABLE").length}</div>
              <div className="text-xs lg:text-sm text-slate-600">Available</div>
            </div>
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-amber-600">{offers.reduce((sum, o) => sum + (o.servings || 0), 0)}</div>
              <div className="text-xs lg:text-sm text-slate-600">Total Servings</div>
            </div>
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-purple-600">{offers.reduce((sum, o) => sum + (o.claimsCount || 0), 0)}</div>
              <div className="text-xs lg:text-sm text-slate-600">Total Claims</div>
            </div>
          </div>
        </div>
      </div>

      {/* Advanced Filters */}
      <div className="modern-card p-4 lg:p-6 animate-slide-in">
        <div className="flex items-center justify-between mb-4 lg:mb-6">
          <h2 className="text-lg lg:text-xl font-bold text-slate-900">Find Perfect Donations</h2>
          <Filter className="w-4 h-4 lg:w-5 lg:h-5 text-emerald-600" />
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Search</label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-400" />
              <input
                type="text"
                placeholder="Search offers, organizations..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Status</label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
            >
              <option value="all">All Status</option>
              <option value="AVAILABLE">Available</option>
              <option value="PENDING_CONFIRMATION">Pending</option>
              <option value="CLAIMED">Claimed</option>
              <option value="EXPIRED">Expired</option>
              <option value="CANCELLED">Cancelled</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Purpose</label>
            <select
              value={purposeFilter}
              onChange={(e) => setPurposeFilter(e.target.value)}
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
            >
              <option value="all">All Purposes</option>
              <option value="GENERAL">General</option>
              <option value="CHILDREN">Children</option>
              <option value="ELDERLY">Elderly</option>
              <option value="WOMEN">Women</option>
              <option value="EMERGENCY">Emergency</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Sort By</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
            >
              <option value="newest">Newest First</option>
              <option value="oldest">Oldest First</option>
              <option value="servings">Most Servings</option>
              <option value="distance">Nearest</option>
              <option value="rating">Highest Rated</option>
            </select>
          </div>
        </div>

        <div className="flex items-center justify-between mt-6 pt-4 border-t border-slate-200">
          <div className="text-sm text-slate-600">
            Showing {sortedOffers.length} of {offers.length} offers
          </div>
          <button 
            onClick={() => {
              setSearchTerm("");
              setStatusFilter("all");
              setPurposeFilter("all");
              setSortBy("newest");
            }}
            className="text-emerald-600 hover:text-emerald-500 font-medium text-sm transition-colors"
          >
            Clear Filters
          </button>
        </div>
      </div>

      {/* Offers Grid */}
      {offersLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="h-48 bg-slate-200 rounded-xl mb-4"></div>
              <div className="h-6 bg-slate-200 rounded w-3/4 mb-3"></div>
              <div className="h-4 bg-slate-200 rounded w-1/2 mb-4"></div>
              <div className="space-y-2">
                <div className="h-3 bg-slate-200 rounded w-full"></div>
                <div className="h-3 bg-slate-200 rounded w-2/3"></div>
              </div>
            </div>
          ))}
        </div>
      ) : sortedOffers.length === 0 ? (
        <div className="modern-card p-16 text-center animate-fade-in">
          <div className="w-24 h-24 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-2xl flex items-center justify-center mx-auto mb-6">
            <Gift className="h-12 w-12 text-white" />
          </div>
          <h3 className="text-2xl font-bold text-slate-900 mb-4">No offers found</h3>
          <p className="text-slate-600 text-lg mb-8 max-w-md mx-auto">
            {searchTerm || statusFilter !== "all" || purposeFilter !== "all"
              ? "Try adjusting your search filters to find more offers"
              : user?.role === "COLLAB"
              ? "Share your first food donation and start making an impact"
              : "No food donations are currently available in your area"
            }
          </p>
          {canCreateOffer && (
            <Link
              href="/dashboard/offers/new"
              className="btn-modern px-8 py-4 text-lg inline-flex items-center group"
            >
              <Plus className="h-6 w-6 mr-3 group-hover:scale-110 transition-transform" />
              Create Your First Offer
            </Link>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
          {sortedOffers.map((offer, index) => {
            const StatusIcon = getStatusIcon(offer.status);
            const PurposeIcon = getPurposeIcon(offer.purpose);
            
            return (
              <div 
                key={offer.id} 
                className="modern-card group hover:scale-105 transition-all duration-300 overflow-hidden animate-slide-in"
                style={{animationDelay: `${index * 0.1}s`}}
              >
                {/* Image Header */}
                <div className="relative h-48 bg-gradient-to-br from-emerald-400 to-cyan-500 overflow-hidden">
                  <div className="absolute inset-0 bg-hero-pattern opacity-20"></div>
                  <div className="absolute top-4 left-4 flex space-x-2">
                    <span className={`px-3 py-1 rounded-full text-xs font-bold bg-gradient-to-r ${getStatusColor(offer.status)} shadow-lg`}>
                      <StatusIcon className="w-3 h-3 inline mr-1" />
                      {offer.status.replace("_", " ")}
                    </span>
                    <span className={`px-3 py-1 rounded-full text-xs font-bold bg-gradient-to-r ${getPriorityColor(offer.priority)} text-white shadow-lg`}>
                      {offer.priority}
                    </span>
                  </div>
                  <div className="absolute top-4 right-4">
                    <div className="flex items-center space-x-1 bg-white/90 backdrop-blur-sm px-2 py-1 rounded-full">
                      <Star className="w-3 h-3 text-amber-500" />
                      <span className="text-xs font-bold text-slate-900">{offer.rating}</span>
                    </div>
                  </div>
                  <div className="absolute bottom-4 left-4 flex space-x-2">
                    {offer.tags.map(tag => (
                      <span key={tag} className="px-2 py-1 bg-white/90 backdrop-blur-sm text-xs font-medium text-slate-700 rounded-lg">
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>

                {/* Content */}
                <div className="p-6">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex-1">
                      <h3 className="text-xl font-bold text-slate-900 mb-2 group-hover:text-emerald-600 transition-colors">
                        {offer.title}
                      </h3>
                      <p className="text-slate-600 text-sm mb-3 line-clamp-2 leading-relaxed">
                        {offer.description}
                      </p>
                      <p className="text-emerald-600 font-semibold text-sm">
                        by {offer.organization}
                      </p>
                    </div>
                  </div>

                  {/* Key Metrics */}
                  <div className="grid grid-cols-2 gap-4 mb-6">
                    <div className="bg-slate-50 p-3 rounded-xl">
                      <div className="flex items-center space-x-2 mb-1">
                        <Users className="w-4 h-4 text-emerald-600" />
                        <span className="text-xs text-slate-500 font-medium">Servings</span>
                      </div>
                      <div className="text-xl font-bold text-slate-900">{offer.servings}</div>
                    </div>
                    <div className="bg-slate-50 p-3 rounded-xl">
                      <div className="flex items-center space-x-2 mb-1">
                        <MapPin className="w-4 h-4 text-emerald-600" />
                        <span className="text-xs text-slate-500 font-medium">Distance</span>
                      </div>
                      <div className="text-xl font-bold text-slate-900">{offer.distance}</div>
                    </div>
                  </div>

                  {/* Details */}
                  <div className="space-y-3 mb-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <PurposeIcon className="w-4 h-4 text-slate-400" />
                        <span className="text-sm text-slate-600">Purpose</span>
                      </div>
                      <span className={`px-2 py-1 rounded-lg text-xs font-bold bg-gradient-to-r ${getPurposeColor(offer.purpose)}`}>
                        {offer.purpose}
                      </span>
                    </div>
                    
                    <div className="flex items-center space-x-2 text-sm text-slate-600">
                      <MapPin className="h-4 w-4 text-slate-400" />
                      <span>{offer.location.area}, {offer.location.city}</span>
                    </div>
                    
                    <div className="flex items-center space-x-2 text-sm text-slate-600">
                      <Calendar className="h-4 w-4 text-slate-400" />
                      <span>Ready: {new Date(offer.readyFrom).toLocaleDateString()}</span>
                    </div>
                    
                    <div className="flex items-center space-x-2 text-sm text-slate-600">
                      <Clock className="h-4 w-4 text-slate-400" />
                      <span>Expires: {new Date(offer.expiresAt).toLocaleDateString()}</span>
                    </div>
                    
                    {offer.claimsCount > 0 && (
                      <div className="flex items-center space-x-2 text-sm text-slate-600">
                        <Activity className="h-4 w-4 text-slate-400" />
                        <span>{offer.claimsCount} organization{offer.claimsCount > 1 ? 's' : ''} interested</span>
                      </div>
                    )}
                  </div>

                  {/* Actions */}
                  <div className="flex space-x-3">
                    <Link
                      href={`/dashboard/offers/${offer.id}`}
                      className="flex-1 btn-modern py-3 text-center group/btn"
                    >
                      {canClaimOffer && offer.status === "AVAILABLE" ? (
                        <>
                          <Heart className="w-4 h-4 mr-2 group-hover/btn:scale-110 transition-transform inline" />
                          Claim Offer
                        </>
                      ) : (
                        <>
                          <Eye className="w-4 h-4 mr-2 group-hover/btn:scale-110 transition-transform inline" />
                          View Details
                        </>
                      )}
                      <ArrowRight className="w-4 h-4 ml-2 group-hover/btn:translate-x-1 transition-transform inline" />
                    </Link>
                    {user?.role === "COLLAB" && offer.createdBy === "COLLAB" && (
                      <button className="px-4 py-3 border-2 border-slate-200 text-slate-600 hover:border-emerald-300 hover:text-emerald-600 hover:bg-emerald-50 rounded-xl text-sm font-semibold transition-all group/edit">
                        <Edit className="h-4 w-4 group-hover/edit:scale-110 transition-transform" />
                      </button>
                    )}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
