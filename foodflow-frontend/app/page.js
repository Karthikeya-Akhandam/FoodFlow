import Link from "next/link";
import { Heart, Users, Shield, Zap, ArrowRight, CheckCircle, Sparkles, Globe, TrendingUp, Star } from "lucide-react";

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-emerald-50 relative overflow-hidden">
      {/* Background decorations */}
      <div className="absolute inset-0 bg-hero-pattern opacity-50"></div>
      <div className="absolute top-20 left-10 w-72 h-72 bg-gradient-to-r from-emerald-400 to-cyan-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-float"></div>
      <div className="absolute top-40 right-10 w-72 h-72 bg-gradient-to-r from-purple-400 to-pink-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-float" style={{animationDelay: '2s'}}></div>
      <div className="absolute bottom-20 left-1/3 w-72 h-72 bg-gradient-to-r from-yellow-400 to-orange-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-float" style={{animationDelay: '4s'}}></div>
      
      {/* Navigation */}
      <nav className="glass sticky top-0 z-50 border-b border-white/20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-20">
            <div className="flex items-center">
              <div className="flex-shrink-0 flex items-center space-x-3">
                <div className="w-10 h-10 bg-gradient-primary rounded-xl flex items-center justify-center">
                  <Sparkles className="w-6 h-6 text-white" />
                </div>
                <h1 className="text-3xl font-bold gradient-text">FoodFlow</h1>
              </div>
            </div>
            <div className="hidden md:block">
              <div className="ml-10 flex items-center space-x-1">
                <Link href="#features" className="text-slate-600 hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all hover:bg-white/50">
                  Features
                </Link>
                <Link href="#how-it-works" className="text-slate-600 hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all hover:bg-white/50">
                  How it Works
                </Link>
                <Link href="/auth/login" className="text-slate-600 hover:text-primary px-4 py-2 rounded-lg text-sm font-medium transition-all hover:bg-white/50">
                  Login
                </Link>
                <Link href="/auth/signup" className="btn-modern ml-4">
                  Get Started
                </Link>
              </div>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative pt-20 pb-32 overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="text-center animate-fade-in">
            <div className="inline-flex items-center px-6 py-3 rounded-full bg-white/90 backdrop-blur-sm border border-emerald-200 text-emerald-700 text-sm font-medium mb-8 hover-lift glow-effect">
              <Star className="w-4 h-4 mr-2 animate-bounce" />
              <span>Connecting communities through food sharing</span>
            </div>
            <h1 className="text-5xl md:text-7xl lg:text-8xl font-bold text-slate-900 mb-8 leading-tight text-shadow">
              Transform Food
              <span className="block gradient-text animate-shimmer">Into Hope</span>
            </h1>
            <p className="text-xl md:text-2xl text-slate-600 mb-12 max-w-4xl mx-auto leading-relaxed">
              Join the revolution in food distribution. Our AI-powered platform connects surplus food 
              with communities in need, reducing waste and fighting hunger one meal at a time.
            </p>
            <div className="flex flex-col sm:flex-row gap-6 justify-center items-center">
              <Link href="/auth/signup" className="btn-modern text-lg px-10 py-4 animate-pulse-glow group hover-lift">
                <span className="flex items-center">
                  Start Your Impact
                  <ArrowRight className="ml-3 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                </span>
              </Link>
              <Link href="/auth/signup?role=ORG" className="btn-secondary text-lg px-10 py-4 rounded-xl hover-lift">
                Join as Organization
              </Link>
            </div>
            <div className="mt-16 flex justify-center items-center space-x-8 text-slate-500">
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-emerald-100 rounded-lg flex items-center justify-center">
                  <Users className="w-4 h-4 text-emerald-600" />
                </div>
                <span className="text-sm font-medium">500+ Organizations</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                  <Globe className="w-4 h-4 text-blue-600" />
                </div>
                <span className="text-sm font-medium">50+ Cities</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-amber-100 rounded-lg flex items-center justify-center">
                  <TrendingUp className="w-4 h-4 text-amber-600" />
                </div>
                <span className="text-sm font-medium">1M+ Meals Shared</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-32 relative">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-20 animate-fade-in">
            <h2 className="text-4xl md:text-6xl font-bold text-slate-900 mb-6">
              Why Choose
              <span className="gradient-text block">FoodFlow?</span>
            </h2>
            <p className="text-xl md:text-2xl text-slate-600 max-w-3xl mx-auto leading-relaxed">
              Revolutionary technology meets compassionate action to create maximum impact
            </p>
          </div>
          
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="modern-card p-8 text-center group animate-slide-in hover-lift glow-effect">
              <div className="w-20 h-20 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform animate-pulse-glow">
                <Heart className="h-10 w-10 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-slate-900 mb-4">Zero Waste</h3>
              <p className="text-slate-600 leading-relaxed">AI-powered matching ensures every meal finds its way to someone who needs it, eliminating food waste</p>
            </div>
            
            <div className="modern-card p-8 text-center group animate-slide-in hover-lift glow-effect" style={{animationDelay: '0.2s'}}>
              <div className="w-20 h-20 bg-gradient-to-br from-blue-400 to-blue-600 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform animate-pulse-glow">
                <Users className="h-10 w-10 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-slate-900 mb-4">Smart Matching</h3>
              <p className="text-slate-600 leading-relaxed">Advanced algorithms connect donors with recipients based on location, timing, and dietary requirements</p>
            </div>
            
            <div className="modern-card p-8 text-center group animate-slide-in hover-lift glow-effect" style={{animationDelay: '0.4s'}}>
              <div className="w-20 h-20 bg-gradient-to-br from-purple-400 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform animate-pulse-glow">
                <Shield className="h-10 w-10 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-slate-900 mb-4">Secure & Verified</h3>
              <p className="text-slate-600 leading-relaxed">Blockchain-verified transactions and comprehensive background checks ensure complete trust and safety</p>
            </div>
            
            <div className="modern-card p-8 text-center group animate-slide-in hover-lift glow-effect" style={{animationDelay: '0.6s'}}>
              <div className="w-20 h-20 bg-gradient-to-br from-amber-400 to-amber-600 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform animate-pulse-glow">
                <Zap className="h-10 w-10 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-slate-900 mb-4">Lightning Fast</h3>
              <p className="text-slate-600 leading-relaxed">Real-time notifications and instant matching get food to those who need it in record time</p>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section id="how-it-works" className="py-32 bg-gradient-to-br from-slate-50 to-emerald-50 relative">
        <div className="absolute inset-0 bg-hero-pattern opacity-30"></div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="text-center mb-20 animate-fade-in">
            <h2 className="text-4xl md:text-6xl font-bold text-slate-900 mb-6">
              How It
              <span className="gradient-text block">Works</span>
            </h2>
            <p className="text-xl md:text-2xl text-slate-600 max-w-3xl mx-auto leading-relaxed">
              Four simple steps to transform surplus food into community impact
            </p>
          </div>
          
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            <div className="space-y-10">
              <div className="flex items-start space-x-6 animate-slide-in group hover-lift">
                <div className="bg-gradient-to-br from-emerald-400 to-emerald-600 text-white rounded-2xl w-16 h-16 flex items-center justify-center font-bold text-xl shadow-lg group-hover:scale-110 transition-transform animate-pulse-glow">1</div>
                <div>
                  <h3 className="text-2xl font-bold text-slate-900 mb-3">Create Your Profile</h3>
                  <p className="text-slate-600 text-lg leading-relaxed">Join our community in seconds. Whether you're a restaurant, grocery store, or organization in need, we'll match you with the perfect partners.</p>
                </div>
              </div>
              
              <div className="flex items-start space-x-6 animate-slide-in group hover-lift" style={{animationDelay: '0.2s'}}>
                <div className="bg-gradient-to-br from-blue-400 to-blue-600 text-white rounded-2xl w-16 h-16 flex items-center justify-center font-bold text-xl shadow-lg group-hover:scale-110 transition-transform animate-pulse-glow">2</div>
                <div>
                  <h3 className="text-2xl font-bold text-slate-900 mb-3">Smart Discovery</h3>
                  <p className="text-slate-600 text-lg leading-relaxed">Our AI instantly matches available food with nearby organizations, considering dietary restrictions, transportation, and timing.</p>
                </div>
              </div>
              
              <div className="flex items-start space-x-6 animate-slide-in group hover-lift" style={{animationDelay: '0.4s'}}>
                <div className="bg-gradient-to-br from-purple-400 to-purple-600 text-white rounded-2xl w-16 h-16 flex items-center justify-center font-bold text-xl shadow-lg group-hover:scale-110 transition-transform animate-pulse-glow">3</div>
                <div>
                  <h3 className="text-2xl font-bold text-slate-900 mb-3">Seamless Coordination</h3>
                  <p className="text-slate-600 text-lg leading-relaxed">Real-time communication tools and logistics support ensure smooth handoffs from donors to recipients.</p>
                </div>
              </div>
              
              <div className="flex items-start space-x-6 animate-slide-in group hover-lift" style={{animationDelay: '0.6s'}}>
                <div className="bg-gradient-to-br from-amber-400 to-amber-600 text-white rounded-2xl w-16 h-16 flex items-center justify-center font-bold text-xl shadow-lg group-hover:scale-110 transition-transform animate-pulse-glow">4</div>
                <div>
                  <h3 className="text-2xl font-bold text-slate-900 mb-3">Track Impact</h3>
                  <p className="text-slate-600 text-lg leading-relaxed">See your positive impact with detailed analytics, community stories, and recognition for your contributions.</p>
                </div>
              </div>
            </div>
            
            <div className="modern-card p-10 animate-fade-in hover-lift glow-effect" style={{animationDelay: '0.8s'}}>
              <h3 className="text-3xl font-bold text-slate-900 mb-8 text-center">Impact Dashboard</h3>
              <div className="space-y-6">
                <div className="flex items-center justify-between p-4 bg-gradient-to-r from-emerald-50 to-emerald-100 rounded-xl">
                  <div className="flex items-center space-x-3">
                    <CheckCircle className="h-6 w-6 text-emerald-600" />
                    <span className="text-slate-700 font-medium">Food Waste Reduced</span>
                  </div>
                  <span className="text-2xl font-bold text-emerald-600">87%</span>
                </div>
                <div className="flex items-center justify-between p-4 bg-gradient-to-r from-blue-50 to-blue-100 rounded-xl">
                  <div className="flex items-center space-x-3">
                    <CheckCircle className="h-6 w-6 text-blue-600" />
                    <span className="text-slate-700 font-medium">Families Fed Monthly</span>
                  </div>
                  <span className="text-2xl font-bold text-blue-600">12K+</span>
                </div>
                <div className="flex items-center justify-between p-4 bg-gradient-to-r from-purple-50 to-purple-100 rounded-xl">
                  <div className="flex items-center space-x-3">
                    <CheckCircle className="h-6 w-6 text-purple-600" />
                    <span className="text-slate-700 font-medium">Community Partners</span>
                  </div>
                  <span className="text-2xl font-bold text-purple-600">500+</span>
                </div>
                <div className="flex items-center justify-between p-4 bg-gradient-to-r from-amber-50 to-amber-100 rounded-xl">
                  <div className="flex items-center space-x-3">
                    <CheckCircle className="h-6 w-6 text-amber-600" />
                    <span className="text-slate-700 font-medium">Carbon Footprint Saved</span>
                  </div>
                  <span className="text-2xl font-bold text-amber-600">2.3M kg</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-32 bg-gradient-primary relative overflow-hidden">
        <div className="absolute inset-0">
          <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-emerald-600/20 to-cyan-600/20"></div>
          <div className="absolute top-20 left-20 w-64 h-64 bg-white/10 rounded-full blur-3xl"></div>
          <div className="absolute bottom-20 right-20 w-96 h-96 bg-white/5 rounded-full blur-3xl"></div>
        </div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
          <div className="animate-fade-in">
            <h2 className="text-4xl md:text-6xl lg:text-7xl font-bold text-white mb-8 leading-tight">
              Ready to Transform
              <span className="block text-emerald-200">Your Community?</span>
            </h2>
            <p className="text-xl md:text-2xl text-emerald-100 mb-12 max-w-4xl mx-auto leading-relaxed">
              Join thousands of changemakers using FoodFlow to eliminate waste, feed communities, 
              and build a more sustainable future—one meal at a time.
            </p>
            <div className="flex flex-col sm:flex-row gap-6 justify-center items-center">
              <Link href="/auth/signup" className="bg-white text-emerald-600 hover:bg-emerald-50 px-10 py-4 rounded-xl text-lg font-bold transition-all transform hover:scale-105 hover:shadow-2xl inline-flex items-center group">
                <span>Start Your Impact Journey</span>
                <ArrowRight className="ml-3 h-5 w-5 group-hover:translate-x-1 transition-transform" />
              </Link>
              <Link href="/auth/signup?role=ORG" className="border-2 border-white/30 text-white hover:bg-white/10 px-10 py-4 rounded-xl text-lg font-bold transition-all backdrop-blur-sm">
                Partner With Us
              </Link>
            </div>
            <div className="mt-16 flex flex-wrap justify-center items-center gap-8 text-emerald-200">
              <div className="text-center">
                <div className="text-3xl font-bold text-white">500K+</div>
                <div className="text-sm opacity-80">Meals Saved</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-white">1,200+</div>
                <div className="text-sm opacity-80">Partners</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-white">50+</div>
                <div className="text-sm opacity-80">Cities</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-bold text-white">95%</div>
                <div className="text-sm opacity-80">Success Rate</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-slate-900 text-white py-16 relative overflow-hidden">
        <div className="absolute inset-0">
          <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-emerald-900/20 to-slate-900"></div>
        </div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="grid md:grid-cols-4 gap-12">
            <div className="md:col-span-2">
              <div className="flex items-center space-x-3 mb-6">
                <div className="w-12 h-12 bg-gradient-primary rounded-xl flex items-center justify-center">
                  <Sparkles className="w-7 h-7 text-white" />
                </div>
                <h3 className="text-3xl font-bold gradient-text">FoodFlow</h3>
              </div>
              <p className="text-slate-400 text-lg leading-relaxed mb-8 max-w-md">
                Revolutionizing food distribution through technology, compassion, and community collaboration. Together, we're building a world without food waste.
              </p>
              <div className="flex space-x-4">
                <div className="w-10 h-10 bg-slate-800 rounded-lg flex items-center justify-center hover:bg-emerald-600 transition-colors cursor-pointer">
                  <Globe className="w-5 h-5" />
                </div>
                <div className="w-10 h-10 bg-slate-800 rounded-lg flex items-center justify-center hover:bg-emerald-600 transition-colors cursor-pointer">
                  <Heart className="w-5 h-5" />
                </div>
                <div className="w-10 h-10 bg-slate-800 rounded-lg flex items-center justify-center hover:bg-emerald-600 transition-colors cursor-pointer">
                  <Users className="w-5 h-5" />
                </div>
              </div>
            </div>
            <div>
              <h4 className="text-xl font-bold mb-6 text-emerald-400">Platform</h4>
              <ul className="space-y-3 text-slate-400">
                <li><Link href="#features" className="hover:text-emerald-400 transition-colors text-lg">Features</Link></li>
                <li><Link href="#how-it-works" className="hover:text-emerald-400 transition-colors text-lg">How it Works</Link></li>
                <li><Link href="/auth/signup" className="hover:text-emerald-400 transition-colors text-lg">Get Started</Link></li>
                <li><Link href="/dashboard" className="hover:text-emerald-400 transition-colors text-lg">Dashboard</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="text-xl font-bold mb-6 text-emerald-400">Support</h4>
              <ul className="space-y-3 text-slate-400">
                <li><Link href="/help" className="hover:text-emerald-400 transition-colors text-lg">Help Center</Link></li>
                <li><Link href="/contact" className="hover:text-emerald-400 transition-colors text-lg">Contact Us</Link></li>
                <li><Link href="/about" className="hover:text-emerald-400 transition-colors text-lg">About Us</Link></li>
                <li><Link href="/blog" className="hover:text-emerald-400 transition-colors text-lg">Blog</Link></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-slate-800 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center">
            <p className="text-slate-400 text-lg">&copy; 2024 FoodFlow. Transforming communities, one meal at a time.</p>
            <div className="mt-4 md:mt-0 text-slate-500">
              <span className="text-sm">Built with ❤️ for a better world</span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}