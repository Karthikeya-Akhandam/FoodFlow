"use client";

import { useMemo } from "react";
import Link from "next/link";
import { 
  Plus, 
  Gift, 
  Users, 
  TrendingUp, 
  Clock, 
  CheckCircle,
  AlertCircle,
  ArrowRight,
  Sparkles,
  Target,
  Calendar,
  MapPin,
  Star,
  Award,
  Activity,
  Heart,
  Zap
} from "lucide-react";
import { useAuth } from "../../lib/auth";
import { useOffers, useClaims, useCredits, useTokens, useProfileCompletion } from "../../lib/hooks";

export default function DashboardPage() {
  const { user } = useAuth();
  const { needsCompletion } = useProfileCompletion();
  
  const { data: offers, isLoading: offersLoading, error: offersError } = useOffers();
  const { data: claims, isLoading: claimsLoading, error: claimsError } = useClaims();
  
  // Only load credits/tokens if user has completed their profile setup
  const { data: credits, isLoading: creditsLoading, error: creditsError } = useCredits({
    enabled: !needsCompletion && user?.role === "ORG"
  });
  const { data: tokens, isLoading: tokensLoading, error: tokensError } = useTokens({
    enabled: !needsCompletion && user?.role === "COLLAB"
  });

  const isLoading = offersLoading || claimsLoading || creditsLoading || tokensLoading;
  const hasError = offersError || claimsError || creditsError || tokensError;

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-8">
        {/* Welcome Section Skeleton */}
        <div className="modern-card p-8 animate-pulse">
          <div className="flex items-center justify-between">
            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <div className="w-12 h-12 bg-slate-200 rounded-xl"></div>
                <div className="space-y-2">
                  <div className="h-8 bg-slate-200 rounded w-64"></div>
                  <div className="h-4 bg-slate-200 rounded w-32"></div>
                </div>
              </div>
              <div className="h-4 bg-slate-200 rounded w-96"></div>
            </div>
            <div className="w-32 h-24 bg-slate-200 rounded"></div>
          </div>
        </div>

        {/* Quick Actions Skeleton */}
        <div className="grid md:grid-cols-2 gap-6">
          {[...Array(2)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="flex items-center space-x-4">
                <div className="w-16 h-16 bg-slate-200 rounded-2xl"></div>
                <div className="space-y-2 flex-1">
                  <div className="h-5 bg-slate-200 rounded w-32"></div>
                  <div className="h-4 bg-slate-200 rounded w-48"></div>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Stats Grid Skeleton */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="modern-card p-6 animate-pulse">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="w-12 h-12 bg-slate-200 rounded-xl"></div>
                  <div className="h-4 bg-slate-200 rounded w-16"></div>
                </div>
                <div className="space-y-2">
                  <div className="h-8 bg-slate-200 rounded w-20"></div>
                  <div className="h-4 bg-slate-200 rounded w-24"></div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  // Error state
  if (hasError) {
    return (
      <div className="space-y-8">
        <div className="text-center py-12">
          <AlertCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-slate-900 mb-2">Unable to Load Dashboard</h2>
          <p className="text-slate-600 mb-6">
            There was an error loading your dashboard data. Please try again later.
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
  
  const stats = useMemo(() => ({
    totalOffers: offers?.length || 0,
    activeOffers: offers?.filter(offer => offer.status === 'AVAILABLE')?.length || 0,
    totalClaims: claims?.length || 0,
    successfulRedemptions: claims?.filter(claim => claim.status === 'WON')?.length || 0,
    credits: credits?.balance || 0,
    tokens: tokens?.balance || 0
  }), [offers, claims, credits, tokens]);

  const getGreeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return "Good morning";
    if (hour < 18) return "Good afternoon";
    return "Good evening";
  };

  const getRoleName = (role) => {
    switch (role) {
      case "ORG": return "Organization";
      case "COLLAB": return "Food Donor";
      case "ADMIN": return "Administrator";
      default: return "User";
    }
  };

  const getRoleDescription = (role) => {
    switch (role) {
      case "ORG": return "Browse and claim food donations from local businesses in your community";
      case "COLLAB": return "Share surplus food and earn tokens for making a positive impact";
      case "ADMIN": return "Manage the platform and oversee all operations to maximize impact";
      default: return "Welcome to FoodFlow - let's make a difference together";
    }
  };

  const getQuickActions = (role) => {
    switch (role) {
      case "ORG":
        return [
          { name: "Browse Offers", href: "/dashboard/offers", icon: Gift, color: "emerald", description: "Find available food donations" },
          { name: "My Claims", href: "/dashboard/claims", icon: Users, color: "blue", description: "Track your claims" },
        ];
      case "COLLAB":
        return [
          { name: "Create Offer", href: "/dashboard/offers/new", icon: Plus, color: "emerald", description: "Share surplus food" },
          { name: "My Offers", href: "/dashboard/offers", icon: Gift, color: "blue", description: "Manage your donations" },
        ];
      default:
        return [];
    }
  };

  const getStatsCards = (role) => {
    if (role === "ORG") {
      return [
        { 
          label: "Total Claims", 
          value: stats.totalClaims, 
          icon: Users, 
          color: "blue",
          change: "+12%",
          trend: "up"
        },
        { 
          label: "Successful Claims", 
          value: stats.successfulRedemptions, 
          icon: CheckCircle, 
          color: "emerald",
          change: "+8%",
          trend: "up"
        },
        { 
          label: "Available Credits", 
          value: stats.credits, 
          icon: Star, 
          color: "amber",
          change: "-5 today",
          trend: "down"
        },
        { 
          label: "Impact Score", 
          value: "95%", 
          icon: Award, 
          color: "purple",
          change: "+2%",
          trend: "up"
        },
      ];
    } else if (role === "COLLAB") {
      return [
        { 
          label: "Total Offers", 
          value: stats.totalOffers, 
          icon: Gift, 
          color: "emerald",
          change: "+3 this week",
          trend: "up"
        },
        { 
          label: "Active Offers", 
          value: stats.activeOffers, 
          icon: Activity, 
          color: "blue",
          change: "2 pending",
          trend: "neutral"
        },
        { 
          label: "Earned Tokens", 
          value: stats.tokens, 
          icon: Zap, 
          color: "amber",
          change: "+15 today",
          trend: "up"
        },
        { 
          label: "Lives Impacted", 
          value: "247", 
          icon: Heart, 
          color: "red",
          change: "+18 this week",
          trend: "up"
        },
      ];
    }
    return [];
  };

  const quickActions = getQuickActions(user?.role);
  const statsCards = getStatsCards(user?.role);

  // Generate recent activity from real data
  const recentActivity = useMemo(() => {
    const activities = [];
    
    // Add recent offers (for COLLAB users)
    if (user?.role === "COLLAB" && offers) {
      offers.slice(0, 2).forEach((offer, index) => {
        activities.push({
          id: `offer-${offer.id}`,
          type: "offer_created",
          title: "You created a new offer",
          description: `${offer.title} - ${offer.servings || 0} servings`,
          time: formatTimeAgo(offer.created_at || offer.createdAt),
          icon: Gift,
          color: "emerald"
        });
      });
    }
    
    // Add recent claims (for ORG users)
    if (user?.role === "ORG" && claims) {
      claims.slice(0, 2).forEach((claim, index) => {
        activities.push({
          id: `claim-${claim.id}`,
          type: "claim_made",
          title: claim.status === "WON" ? "Claim won!" : "Claim submitted",
          description: `${claim.offer_title || 'Food offer'} - ${claim.requested_servings || 0} servings`,
          time: formatTimeAgo(claim.created_at || claim.createdAt),
          icon: claim.status === "WON" ? CheckCircle : Clock,
          color: claim.status === "WON" ? "emerald" : "blue"
        });
      });
    }
    
    // Add fallback activities if no real data
    if (activities.length === 0) {
      activities.push({
        id: 1,
        type: "welcome",
        title: "Welcome to FoodFlow!",
        description: "Start by exploring available offers or creating your first donation",
        time: "Now",
        icon: Sparkles,
        color: "purple"
      });
    }
    
    return activities.slice(0, 3); // Limit to 3 activities
  }, [offers, claims, user?.role]);

  // Helper function to format time ago
  const formatTimeAgo = (dateString) => {
    if (!dateString) return "Recently";
    const date = new Date(dateString);
    const now = new Date();
    const diffInMinutes = Math.floor((now - date) / (1000 * 60));
    
    if (diffInMinutes < 60) return `${diffInMinutes} minutes ago`;
    if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)} hours ago`;
    return `${Math.floor(diffInMinutes / 1440)} days ago`;
  };

  return (
    <div className="space-y-6 lg:space-y-8">
      {/* Welcome Section */}
      <div className="modern-card p-6 lg:p-8 bg-gradient-to-br from-emerald-50 to-cyan-50 border-emerald-200 animate-fade-in">
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
          <div className="flex-1">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-10 h-10 lg:w-12 lg:h-12 bg-gradient-primary rounded-xl flex items-center justify-center flex-shrink-0">
                <Sparkles className="w-5 h-5 lg:w-6 lg:h-6 text-white" />
              </div>
              <div className="min-w-0">
                <h1 className="text-2xl lg:text-3xl font-bold text-slate-900 truncate">
                  {getGreeting()}, {user?.profile?.first_name || "User"}!
                </h1>
                <p className="text-emerald-600 font-medium text-sm lg:text-base">
                  {getRoleName(user?.role)} Dashboard
                </p>
              </div>
            </div>
            <p className="text-slate-600 text-base lg:text-lg leading-relaxed">
              {getRoleDescription(user?.role)}
            </p>
          </div>
          <div className="flex-shrink-0 animate-slide-in">
            <div className="modern-card p-4 bg-white/80 min-w-0">
              <div className="flex items-center space-x-2 mb-2">
                <Calendar className="w-4 h-4 lg:w-5 lg:h-5 text-emerald-600 flex-shrink-0" />
                <span className="text-xs lg:text-sm text-slate-500 font-medium">Today</span>
              </div>
              <p className="text-lg lg:text-xl font-bold text-slate-900 leading-tight">
                {new Date().toLocaleDateString('en-US', { 
                  weekday: 'long',
                  month: 'short',
                  day: 'numeric'
                })}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 lg:gap-6 animate-slide-in">
        {quickActions.map((action, index) => {
          const Icon = action.icon;
          return (
            <Link
              key={action.name}
              href={action.href}
              className="modern-card p-4 lg:p-6 hover:scale-105 transition-all group"
              style={{animationDelay: `${index * 0.1}s`}}
            >
              <div className="flex items-center space-x-3 lg:space-x-4">
                <div className={`w-12 h-12 lg:w-16 lg:h-16 bg-gradient-to-br from-${action.color}-400 to-${action.color}-600 rounded-xl lg:rounded-2xl flex items-center justify-center group-hover:scale-110 transition-transform flex-shrink-0`}>
                  <Icon className="w-6 h-6 lg:w-8 lg:h-8 text-white" />
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="text-lg lg:text-xl font-bold text-slate-900 mb-1 truncate">{action.name}</h3>
                  <p className="text-slate-600 text-sm lg:text-base leading-snug">{action.description}</p>
                </div>
                <ArrowRight className="w-4 h-4 lg:w-5 lg:h-5 text-slate-400 group-hover:text-emerald-600 group-hover:translate-x-1 transition-all flex-shrink-0" />
              </div>
            </Link>
          );
        })}
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 lg:gap-6 animate-fade-in">
        {statsCards.map((stat, index) => {
          const Icon = stat.icon;
          return (
            <div 
              key={stat.label} 
              className="modern-card p-4 lg:p-6 group hover:scale-105 transition-all"
              style={{animationDelay: `${index * 0.1}s`}}
            >
              <div className="flex items-center justify-between mb-3 lg:mb-4">
                <div className={`w-10 h-10 lg:w-12 lg:h-12 bg-gradient-to-br from-${stat.color}-400 to-${stat.color}-600 rounded-lg lg:rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform`}>
                  <Icon className="w-5 h-5 lg:w-6 lg:h-6 text-white" />
                </div>
                <div className={`text-xs lg:text-sm font-medium truncate ml-2 ${stat.trend === 'up' ? 'text-emerald-600' : stat.trend === 'down' ? 'text-red-600' : 'text-slate-500'}`}>
                  {stat.change}
                </div>
              </div>
              <div>
                <p className="text-2xl lg:text-3xl font-bold text-slate-900 mb-1 leading-none">{stat.value}</p>
                <p className="text-slate-600 font-medium text-sm lg:text-base leading-tight">{stat.label}</p>
              </div>
            </div>
          );
        })}
      </div>

      {/* Recent Activity & Impact */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">
        {/* Recent Activity */}
        <div className="lg:col-span-2">
          <div className="modern-card p-4 lg:p-6 animate-slide-in">
            <div className="flex items-center justify-between mb-4 lg:mb-6">
              <h2 className="text-xl lg:text-2xl font-bold text-slate-900">Recent Activity</h2>
              <Link href="/dashboard/activity" className="text-emerald-600 hover:text-emerald-500 font-medium flex items-center text-sm lg:text-base">
                View All
                <ArrowRight className="w-3 h-3 lg:w-4 lg:h-4 ml-1" />
              </Link>
            </div>
            <div className="space-y-3 lg:space-y-4">
              {recentActivity.map((activity, index) => {
                const Icon = activity.icon;
                return (
                  <div 
                    key={activity.id} 
                    className="flex items-center space-x-3 lg:space-x-4 p-3 lg:p-4 rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors"
                    style={{animationDelay: `${index * 0.1}s`}}
                  >
                    <div className={`w-8 h-8 lg:w-10 lg:h-10 bg-gradient-to-br from-${activity.color}-400 to-${activity.color}-600 rounded-lg flex items-center justify-center flex-shrink-0`}>
                      <Icon className="w-4 h-4 lg:w-5 lg:h-5 text-white" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-slate-900 text-sm lg:text-base truncate">{activity.title}</h3>
                      <p className="text-slate-600 text-xs lg:text-sm leading-snug">{activity.description}</p>
                    </div>
                    <div className="text-xs text-slate-500 flex-shrink-0">
                      {activity.time}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Impact Summary */}
        <div className="space-y-4 lg:space-y-6">
          <div className="modern-card p-4 lg:p-6 bg-gradient-to-br from-purple-50 to-pink-50 border-purple-200 animate-fade-in">
            <h2 className="text-lg lg:text-xl font-bold text-slate-900 mb-3 lg:mb-4">Your Impact</h2>
            <div className="space-y-3 lg:space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-slate-600 text-sm lg:text-base">Meals Provided</span>
                <span className="text-xl lg:text-2xl font-bold text-purple-600">1,247</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-slate-600 text-sm lg:text-base">Food Waste Saved</span>
                <span className="text-xl lg:text-2xl font-bold text-purple-600">892 kg</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-slate-600 text-sm lg:text-base">Community Reach</span>
                <span className="text-xl lg:text-2xl font-bold text-purple-600">15 areas</span>
              </div>
            </div>
            <div className="mt-4 lg:mt-6 p-3 lg:p-4 bg-white/80 rounded-xl">
              <p className="text-xs lg:text-sm text-slate-600 text-center leading-snug">
                <span className="font-bold text-purple-600">Amazing work!</span> You're in the top 10% of contributors this month.
              </p>
            </div>
          </div>

          <div className="modern-card p-4 lg:p-6 animate-slide-in">
            <h2 className="text-lg lg:text-xl font-bold text-slate-900 mb-3 lg:mb-4">Quick Links</h2>
            <div className="space-y-2 lg:space-y-3">
              <Link href="/dashboard/profile" className="flex items-center space-x-3 p-2 lg:p-3 rounded-lg hover:bg-slate-50 transition-colors group">
                <Target className="w-4 h-4 lg:w-5 lg:h-5 text-slate-400 group-hover:text-emerald-600 flex-shrink-0" />
                <span className="text-slate-600 group-hover:text-slate-900 text-sm lg:text-base">Update Profile</span>
              </Link>
              <Link href="/dashboard/settings" className="flex items-center space-x-3 p-2 lg:p-3 rounded-lg hover:bg-slate-50 transition-colors group">
                <MapPin className="w-4 h-4 lg:w-5 lg:h-5 text-slate-400 group-hover:text-emerald-600 flex-shrink-0" />
                <span className="text-slate-600 group-hover:text-slate-900 text-sm lg:text-base">Location Settings</span>
              </Link>
              <Link href="/help" className="flex items-center space-x-3 p-2 lg:p-3 rounded-lg hover:bg-slate-50 transition-colors group">
                <AlertCircle className="w-4 h-4 lg:w-5 lg:h-5 text-slate-400 group-hover:text-emerald-600 flex-shrink-0" />
                <span className="text-slate-600 group-hover:text-slate-900 text-sm lg:text-base">Help & Support</span>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}