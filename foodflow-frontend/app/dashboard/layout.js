"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { 
  Home, 
  Gift, 
  Users, 
  User, 
  Settings, 
  LogOut, 
  Menu, 
  X,
  Bell,
  CreditCard,
  TrendingUp,
  Sparkles,
  Search,
  Plus
} from "lucide-react";
import { useAuth, withAuth } from "../../lib/auth";
import { useCredits, useTokens } from "../../lib/hooks";

function DashboardLayout({ children }) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user, logout } = useAuth();
  const pathname = usePathname();
  
  // Fetch credits/tokens based on user role
  const { data: credits } = useCredits();
  const { data: tokens } = useTokens();

  const navigation = [
    { name: "Dashboard", href: "/dashboard", icon: Home },
    { name: "Offers", href: "/dashboard/offers", icon: Gift },
    { name: "Claims", href: "/dashboard/claims", icon: Users },
    { name: "Profile", href: "/dashboard/profile", icon: User },
    { name: "Settings", href: "/dashboard/settings", icon: Settings },
  ];

  const isActive = (href) => {
    if (href === "/dashboard") {
      return pathname === "/dashboard";
    }
    return pathname.startsWith(href);
  };

  const handleLogout = () => {
    logout();
    window.location.href = '/';
  };

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-emerald-50 relative overflow-hidden">
      {/* Background decorations */}
      <div className="absolute inset-0 bg-hero-pattern opacity-20"></div>
      <div className="absolute top-20 right-20 w-96 h-96 bg-gradient-to-r from-emerald-400/20 to-cyan-400/20 rounded-full mix-blend-multiply filter blur-3xl"></div>
      <div className="absolute bottom-20 left-20 w-96 h-96 bg-gradient-to-r from-purple-400/20 to-pink-400/20 rounded-full mix-blend-multiply filter blur-3xl"></div>
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? "block" : "hidden"}`}>
        <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm" onClick={() => setSidebarOpen(false)} />
        <div className="fixed inset-y-0 left-0 flex w-80 flex-col glass border-r border-white/20">
          <div className="flex h-20 items-center justify-between px-6">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-primary rounded-xl flex items-center justify-center">
                <Sparkles className="w-6 h-6 text-white" />
              </div>
              <h1 className="text-2xl font-bold gradient-text">FoodFlow</h1>
            </div>
            <button
              onClick={() => setSidebarOpen(false)}
              className="text-slate-500 hover:text-slate-700 transition-colors"
            >
              <X className="h-6 w-6" />
            </button>
          </div>
          <nav className="flex-1 space-y-2 px-4 py-6">
            {navigation.map((item) => {
              const Icon = item.icon;
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`group flex items-center px-4 py-3 text-sm font-semibold rounded-xl transition-all ${
                    isActive(item.href)
                      ? "bg-gradient-primary text-white shadow-lg"
                      : "text-slate-600 hover:bg-white/50 hover:text-emerald-600"
                  }`}
                  onClick={() => setSidebarOpen(false)}
                >
                  <Icon className="mr-3 h-5 w-5" />
                  {item.name}
                </Link>
              );
            })}
          </nav>
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-80 lg:flex-col lg:z-30">
        <div className="flex flex-col flex-grow glass border-r border-white/20">
          <div className="flex h-20 items-center px-6">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-primary rounded-xl flex items-center justify-center">
                <Sparkles className="w-6 h-6 text-white" />
              </div>
              <h1 className="text-2xl font-bold gradient-text">FoodFlow</h1>
            </div>
          </div>
          <nav className="flex-1 space-y-2 px-4 py-6">
            {navigation.map((item) => {
              const Icon = item.icon;
              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`group flex items-center px-4 py-3 text-sm font-semibold rounded-xl transition-all transform hover:scale-105 ${
                    isActive(item.href)
                      ? "bg-gradient-primary text-white shadow-lg"
                      : "text-slate-600 hover:bg-white/50 hover:text-emerald-600"
                  }`}
                >
                  <Icon className="mr-3 h-5 w-5" />
                  {item.name}
                </Link>
              );
            })}
          </nav>
          
          {/* User info */}
          <div className="border-t border-white/20 p-6 mt-auto">
            <div className="modern-card p-4">
              <div className="flex items-center mb-4">
                <div className="flex-shrink-0">
                  <div className="h-12 w-12 rounded-xl bg-gradient-primary flex items-center justify-center">
                    <span className="text-white text-lg font-bold">
                      {user?.profile?.first_name?.charAt(0) || user?.email?.charAt(0)}
                    </span>
                  </div>
                </div>
                <div className="ml-4">
                  <p className="text-sm font-bold text-slate-900">{user?.profile?.first_name} {user?.profile?.last_name}</p>
                  <p className="text-xs text-emerald-600 font-medium">{user?.role}</p>
                </div>
              </div>
              <button
                onClick={handleLogout}
                className="w-full flex items-center justify-center px-4 py-2 text-sm font-semibold text-slate-600 hover:text-red-600 rounded-lg transition-colors group"
              >
                <LogOut className="mr-2 h-4 w-4 group-hover:scale-110 transition-transform" />
                Sign out
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-80 relative z-10">
        {/* Top bar */}
        <div className="sticky top-0 z-40 glass border-b border-white/20">
          <div className="flex h-20 items-center justify-between px-4 sm:px-6 lg:px-8">
            <div className="flex items-center space-x-4">
              <button
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden text-slate-500 hover:text-slate-700 transition-colors"
              >
                <Menu className="h-6 w-6" />
              </button>
              
              {/* Search bar */}
              <div className="hidden md:block">
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Search className="h-5 w-5 text-slate-400" />
                  </div>
                  <input
                    type="text"
                    placeholder="Search..."
                    className="w-80 pl-10 pr-4 py-2 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  />
                </div>
              </div>
            </div>
            
            <div className="flex items-center space-x-6">
              {/* Credits/Tokens display */}
              {user && (
                <div className="hidden sm:flex items-center space-x-4">
                  {user.role === "ORG" && credits && (
                    <div className="modern-card px-4 py-2 flex items-center space-x-2">
                      <div className="w-8 h-8 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-lg flex items-center justify-center">
                        <CreditCard className="h-4 w-4 text-white" />
                      </div>
                      <div>
                        <div className="text-lg font-bold text-slate-900">{credits.balance || 0}</div>
                        <div className="text-xs text-slate-500 -mt-1">Credits</div>
                      </div>
                    </div>
                  )}
                  {user.role === "COLLAB" && tokens && (
                    <div className="modern-card px-4 py-2 flex items-center space-x-2">
                      <div className="w-8 h-8 bg-gradient-to-br from-amber-400 to-amber-600 rounded-lg flex items-center justify-center">
                        <TrendingUp className="h-4 w-4 text-white" />
                      </div>
                      <div>
                        <div className="text-lg font-bold text-slate-900">{tokens.balance || 0}</div>
                        <div className="text-xs text-slate-500 -mt-1">Tokens</div>
                      </div>
                    </div>
                  )}
                </div>
              )}
              
              {/* Quick Actions */}
              <button className="hidden md:flex items-center space-x-2 btn-modern px-4 py-2">
                <Plus className="h-4 w-4" />
                <span>New</span>
              </button>
              
              {/* Notifications */}
              <button className="relative p-2 text-slate-500 hover:text-emerald-600 transition-colors">
                <Bell className="h-6 w-6" />
                <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
              </button>
              
              {/* User menu */}
              <div className="flex items-center space-x-3">
                <div className="hidden sm:block text-right">
                  <p className="text-sm font-bold text-slate-900">{user?.profile?.first_name} {user?.profile?.last_name}</p>
                  <p className="text-xs text-emerald-600 font-medium">{user?.email}</p>
                </div>
                <div className="h-10 w-10 rounded-xl bg-gradient-primary flex items-center justify-center shadow-lg">
                  <span className="text-white text-sm font-bold">
                    {user?.profile?.first_name?.charAt(0) || user?.email?.charAt(0)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="py-6 lg:py-8">
          <div className="mx-auto max-w-7xl px-3 sm:px-4 lg:px-6 xl:px-8">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}

// Export with authentication protection
export default withAuth(DashboardLayout);
