"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { 
  Search, 
  Filter, 
  MapPin, 
  Clock, 
  Users, 
  MoreVertical,
  Eye,
  CheckCircle,
  XCircle,
  AlertCircle,
  ArrowRight,
  Gift,
  Star,
  TrendingUp,
  Calendar,
  Target,
  Sparkles,
  Activity,
  Heart,
  Zap,
  Award,
  Trophy,
  Shield,
  Plus
} from "lucide-react";
import { useAuth } from "../../../lib/auth";
import { useClaims } from "../../../lib/hooks";

export default function ClaimsPage() {
  const { user } = useAuth();
  const { data: apiClaims, isLoading: claimsLoading, error: claimsError } = useClaims();
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [sortBy, setSortBy] = useState("newest");

  // Use real API data instead of mock data
  const claims = apiClaims || [];

  // Loading state
  if (claimsLoading) {
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

        {/* Stats Cards Skeleton */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="h-4 bg-slate-200 rounded w-3/4 mb-4"></div>
              <div className="h-8 bg-slate-200 rounded w-1/2"></div>
            </div>
          ))}
        </div>

        {/* Claims List Skeleton */}
        <div className="space-y-4">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="h-4 bg-slate-200 rounded w-3/4 mb-2"></div>
              <div className="h-3 bg-slate-200 rounded w-full mb-4"></div>
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
  if (claimsError) {
    return (
      <div className="space-y-8">
        <div className="text-center py-12">
          <AlertCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-slate-900 mb-2">Unable to Load Claims</h2>
          <p className="text-slate-600 mb-6">
            {claimsError.message || "There was an error loading your claims. Please try again later."}
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
      case "REQUESTED": return "from-amber-400 to-amber-600 text-white";
      case "WON": return "from-emerald-400 to-emerald-600 text-white";
      case "LOST": return "from-red-400 to-red-600 text-white";
      case "CANCELLED": return "from-slate-400 to-slate-600 text-white";
      default: return "from-slate-400 to-slate-600 text-white";
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case "REQUESTED": return Clock;
      case "WON": return CheckCircle;
      case "LOST": return XCircle;
      case "CANCELLED": return AlertCircle;
      default: return Clock;
    }
  };

  const getImpactColor = (impact) => {
    switch (impact) {
      case "Critical": return "from-red-400 to-red-600";
      case "High": return "from-orange-400 to-orange-600";
      case "Medium": return "from-yellow-400 to-yellow-600";
      case "Low": return "from-green-400 to-green-600";
      default: return "from-slate-400 to-slate-600";
    }
  };

  const getWinChanceColor = (chance) => {
    if (chance >= 80) return "text-emerald-600";
    if (chance >= 60) return "text-amber-600";
    if (chance >= 40) return "text-orange-600";
    return "text-red-600";
  };

  const filteredClaims = claims.filter(claim => {
    const matchesSearch = claim.offerTitle.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         claim.offerDescription.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         claim.collaborator.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         claim.organization.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === "all" || claim.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const sortedClaims = [...filteredClaims].sort((a, b) => {
    switch (sortBy) {
      case "newest": return new Date(b.createdAt) - new Date(a.createdAt);
      case "oldest": return new Date(a.createdAt) - new Date(b.createdAt);
      case "priority": return b.priorityScore - a.priorityScore;
      case "servings": return b.requestedServings - a.requestedServings;
      case "value": return b.estimatedValue - a.estimatedValue;
      default: return new Date(b.createdAt) - new Date(a.createdAt);
    }
  });

  const getStatusStats = () => {
    const stats = {
      total: claims.length,
      requested: claims.filter(c => c.status === "REQUESTED").length,
      won: claims.filter(c => c.status === "WON").length,
      lost: claims.filter(c => c.status === "LOST").length,
    };
    return stats;
  };

  const stats = getStatusStats();

  return (
    <div className="space-y-6 lg:space-y-8">
      {/* Header Section */}
      <div className="modern-card p-6 lg:p-8 bg-gradient-to-br from-blue-50 to-purple-50 border-blue-200 animate-fade-in">
        <div className="flex flex-col gap-6">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
            <div className="flex items-center space-x-3 lg:space-x-4">
              <div className="w-12 h-12 lg:w-16 lg:h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl lg:rounded-2xl flex items-center justify-center flex-shrink-0">
                <Trophy className="w-6 h-6 lg:w-8 lg:h-8 text-white" />
              </div>
              <div className="min-w-0">
                <h1 className="text-2xl lg:text-4xl font-bold text-slate-900 leading-tight">My Claims Dashboard</h1>
                <p className="text-blue-600 font-medium text-sm lg:text-lg">
                  Track your food donation claims and success rate
                </p>
              </div>
            </div>
            <div className="flex-shrink-0">
              <Link
                href="/dashboard/offers"
                className="btn-modern px-6 lg:px-8 py-3 lg:py-4 text-base lg:text-lg inline-flex items-center group w-full lg:w-auto justify-center"
              >
                <Plus className="h-5 w-5 lg:h-6 lg:w-6 mr-2 lg:mr-3 group-hover:scale-110 transition-transform" />
                Browse New Offers
              </Link>
            </div>
          </div>
          
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-slate-900">{stats.total}</div>
              <div className="text-xs lg:text-sm text-slate-600">Total Claims</div>
            </div>
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-emerald-600">{stats.won}</div>
              <div className="text-xs lg:text-sm text-slate-600">Won</div>
            </div>
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-amber-600">{stats.requested}</div>
              <div className="text-xs lg:text-sm text-slate-600">Pending</div>
            </div>
            <div className="bg-white/80 p-3 lg:p-4 rounded-xl">
              <div className="text-xl lg:text-2xl font-bold text-purple-600">
                {stats.total > 0 ? Math.round((stats.won / stats.total) * 100) : 0}%
              </div>
              <div className="text-xs lg:text-sm text-slate-600">Win Rate</div>
            </div>
          </div>
        </div>
      </div>

      {/* Advanced Filters */}
      <div className="modern-card p-4 lg:p-6 animate-slide-in">
        <div className="flex items-center justify-between mb-4 lg:mb-6">
          <h2 className="text-lg lg:text-xl font-bold text-slate-900">Filter & Sort Claims</h2>
          <Filter className="w-4 h-4 lg:w-5 lg:h-5 text-blue-600" />
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Search</label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-400" />
              <input
                type="text"
                placeholder="Search claims, organizations..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Status</label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
            >
              <option value="all">All Status</option>
              <option value="REQUESTED">Pending</option>
              <option value="WON">Won</option>
              <option value="LOST">Lost</option>
              <option value="CANCELLED">Cancelled</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-semibold text-slate-700 mb-2">Sort By</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="w-full px-4 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
            >
              <option value="newest">Newest First</option>
              <option value="oldest">Oldest First</option>
              <option value="priority">Highest Priority</option>
              <option value="servings">Most Servings</option>
              <option value="value">Highest Value</option>
            </select>
          </div>
        </div>

        <div className="flex items-center justify-between mt-6 pt-4 border-t border-slate-200">
          <div className="text-sm text-slate-600">
            Showing {sortedClaims.length} of {claims.length} claims
          </div>
          <button 
            onClick={() => {
              setSearchTerm("");
              setStatusFilter("all");
              setSortBy("newest");
            }}
            className="text-blue-600 hover:text-blue-500 font-medium text-sm transition-colors"
          >
            Clear Filters
          </button>
        </div>
      </div>

      {/* Claims List */}
      {claimsLoading ? (
        <div className="space-y-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="flex justify-between items-start mb-4">
                <div className="flex-1">
                  <div className="h-6 bg-slate-200 rounded w-3/4 mb-3"></div>
                  <div className="h-4 bg-slate-200 rounded w-1/2 mb-4"></div>
                </div>
                <div className="h-8 w-20 bg-slate-200 rounded-full"></div>
              </div>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                <div className="h-16 bg-slate-200 rounded-xl"></div>
                <div className="h-16 bg-slate-200 rounded-xl"></div>
                <div className="h-16 bg-slate-200 rounded-xl"></div>
                <div className="h-16 bg-slate-200 rounded-xl"></div>
              </div>
              <div className="h-12 bg-slate-200 rounded-xl"></div>
            </div>
          ))}
        </div>
      ) : sortedClaims.length === 0 ? (
        <div className="modern-card p-16 text-center animate-fade-in">
          <div className="w-24 h-24 bg-gradient-to-br from-blue-400 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-6">
            <Users className="h-12 w-12 text-white" />
          </div>
          <h3 className="text-2xl font-bold text-slate-900 mb-4">No claims found</h3>
          <p className="text-slate-600 text-lg mb-8 max-w-md mx-auto">
            {searchTerm || statusFilter !== "all"
              ? "Try adjusting your search filters to find more claims"
              : "Start claiming food donations to help your community"
            }
          </p>
          <Link
            href="/dashboard/offers"
            className="btn-modern px-8 py-4 text-lg inline-flex items-center group"
          >
            <Gift className="h-6 w-6 mr-3 group-hover:scale-110 transition-transform" />
            Browse Available Offers
          </Link>
        </div>
      ) : (
        <div className="space-y-6">
          {sortedClaims.map((claim, index) => {
            const StatusIcon = getStatusIcon(claim.status);
            
            return (
              <div 
                key={claim.id} 
                className="modern-card group hover:scale-102 transition-all duration-300 animate-slide-in"
                style={{animationDelay: `${index * 0.1}s`}}
              >
                <div className="p-8">
                  {/* Header */}
                  <div className="flex items-start justify-between mb-6">
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-3">
                        <h3 className="text-2xl font-bold text-slate-900 group-hover:text-blue-600 transition-colors">
                          {claim.offerTitle}
                        </h3>
                        {claim.isUrgent && (
                          <span className="px-2 py-1 bg-gradient-to-r from-red-400 to-red-600 text-white text-xs font-bold rounded-lg animate-pulse">
                            URGENT
                          </span>
                        )}
                      </div>
                      <p className="text-slate-600 mb-3 leading-relaxed">
                        {claim.offerDescription}
                      </p>
                      <div className="flex items-center space-x-4 text-sm">
                        <span className="text-blue-600 font-semibold">
                          by {claim.collaborator}
                        </span>
                        <span className="text-slate-500">•</span>
                        <span className="text-emerald-600 font-semibold">
                          for {claim.organization}
                        </span>
                        {claim.onBehalfOfRemoteOrg && (
                          <>
                            <span className="text-slate-500">•</span>
                            <span className="text-purple-600 font-semibold">
                              Remote: {claim.onBehalfOfRemoteOrg.name}
                            </span>
                          </>
                        )}
                      </div>
                    </div>
                    <div className="flex items-center space-x-3">
                      <span className={`px-4 py-2 rounded-xl text-sm font-bold bg-gradient-to-r ${getStatusColor(claim.status)} shadow-lg flex items-center space-x-2`}>
                        <StatusIcon className="w-4 h-4" />
                        <span>{claim.status}</span>
                      </span>
                    </div>
                  </div>

                  {/* Key Metrics Grid */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                    <div className="bg-slate-50 p-4 rounded-xl">
                      <div className="flex items-center space-x-2 mb-2">
                        <Users className="w-5 h-5 text-blue-600" />
                        <span className="text-xs text-slate-500 font-medium">Servings</span>
                      </div>
                      <div className="text-2xl font-bold text-slate-900">{claim.requestedServings}</div>
                    </div>
                    
                    <div className="bg-slate-50 p-4 rounded-xl">
                      <div className="flex items-center space-x-2 mb-2">
                        <Star className="w-5 h-5 text-amber-600" />
                        <span className="text-xs text-slate-500 font-medium">Priority Score</span>
                      </div>
                      <div className="text-2xl font-bold text-slate-900">{claim.priorityScore}</div>
                    </div>
                    
                    <div className="bg-slate-50 p-4 rounded-xl">
                      <div className="flex items-center space-x-2 mb-2">
                        <TrendingUp className="w-5 h-5 text-emerald-600" />
                        <span className="text-xs text-slate-500 font-medium">Win Chance</span>
                      </div>
                      <div className={`text-2xl font-bold ${getWinChanceColor(claim.winChance)}`}>
                        {claim.winChance}%
                      </div>
                    </div>
                    
                    <div className="bg-slate-50 p-4 rounded-xl">
                      <div className="flex items-center space-x-2 mb-2">
                        <Activity className="w-5 h-5 text-purple-600" />
                        <span className="text-xs text-slate-500 font-medium">Competitors</span>
                      </div>
                      <div className="text-2xl font-bold text-slate-900">{claim.competitorCount}</div>
                    </div>
                  </div>

                  {/* Details Section */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                    <div className="space-y-3">
                      <div className="flex items-center space-x-3 text-slate-600">
                        <MapPin className="w-4 h-4 text-slate-400" />
                        <span>{claim.offerLocation.area}, {claim.offerLocation.city}</span>
                      </div>
                      
                      <div className="flex items-center space-x-3 text-slate-600">
                        <Calendar className="w-4 h-4 text-slate-400" />
                        <span>Claimed: {new Date(claim.createdAt).toLocaleDateString()}</span>
                      </div>
                      
                      <div className="flex items-center space-x-3 text-slate-600">
                        <Clock className="w-4 h-4 text-slate-400" />
                        <span>Expires: {new Date(claim.offerExpiresAt).toLocaleDateString()}</span>
                      </div>
                    </div>
                    
                    <div className="space-y-3">
                      <div className="flex items-center justify-between">
                        <span className="text-slate-600">Estimated Value</span>
                        <span className="font-bold text-slate-900">₹{claim.estimatedValue.toLocaleString()}</span>
                      </div>
                      
                      <div className="flex items-center justify-between">
                        <span className="text-slate-600">Impact Level</span>
                        <span className={`px-3 py-1 rounded-lg text-xs font-bold bg-gradient-to-r ${getImpactColor(claim.impactScore)} text-white`}>
                          {claim.impactScore}
                        </span>
                      </div>
                      
                      <div className="flex items-center justify-between">
                        <span className="text-slate-600">Category</span>
                        <span className="font-semibold text-slate-900">{claim.category}</span>
                      </div>
                    </div>
                  </div>

                  {/* Tags */}
                  <div className="flex flex-wrap gap-2 mb-6">
                    {claim.tags.map(tag => (
                      <span key={tag} className="px-3 py-1 bg-blue-100 text-blue-700 text-sm font-medium rounded-lg">
                        {tag}
                      </span>
                    ))}
                  </div>

                  {/* Action Buttons */}
                  <div className="flex flex-wrap gap-3">
                    <Link
                      href={`/dashboard/claims/${claim.id}`}
                      className="flex-1 min-w-[200px] btn-modern py-3 text-center group/btn"
                    >
                      <Eye className="w-5 h-5 mr-2 group-hover/btn:scale-110 transition-transform inline" />
                      View Full Details
                      <ArrowRight className="w-4 h-4 ml-2 group-hover/btn:translate-x-1 transition-transform inline" />
                    </Link>
                    
                    {claim.status === "WON" && claim.redemptionStatus === "PENDING_PICKUP" && (
                      <Link
                        href={`/dashboard/redemptions/new?claimId=${claim.id}`}
                        className="btn-modern bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-600 hover:to-emerald-700 py-3 px-6 group/pickup"
                      >
                        <CheckCircle className="w-5 h-5 mr-2 group-hover/pickup:scale-110 transition-transform inline" />
                        Schedule Pickup
                      </Link>
                    )}
                    
                    {claim.status === "WON" && claim.redemptionStatus === "COMPLETED" && (
                      <div className="flex items-center px-6 py-3 bg-emerald-100 text-emerald-700 rounded-xl font-semibold">
                        <CheckCircle className="w-5 h-5 mr-2" />
                        Completed
                      </div>
                    )}
                    
                    {claim.status === "REQUESTED" && (
                      <button className="px-6 py-3 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white rounded-xl font-semibold transition-all group/cancel">
                        <XCircle className="w-5 h-5 mr-2 group-hover/cancel:scale-110 transition-transform inline" />
                        Cancel Claim
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
