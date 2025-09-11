"use client";

import { useState, useEffect } from "react";
import { 
  User, 
  MapPin, 
  Phone, 
  Mail, 
  Building, 
  Users, 
  Edit,
  Save,
  X,
  CheckCircle,
  Sparkles,
  Star,
  Award,
  Heart,
  Target,
  Shield,
  Camera,
  Upload,
  Eye,
  Calendar,
  Activity,
  TrendingUp,
  Clock,
  Zap,
  Globe
} from "lucide-react";
import { useAuth } from "../../../lib/auth";
import { useProfile, useOrganizations, useCollaborators, useUpdateProfile, useUpdateOrganization, useUpdateCollaborator } from "../../../lib/hooks";

export default function ProfilePage() {
  const { user: authUser } = useAuth();
  const { data: profile, isLoading: profileLoading, error: profileError } = useProfile();
  const { data: organizations, isLoading: orgsLoading, error: orgsError } = useOrganizations();
  const { data: collaborators, isLoading: collabsLoading, error: collabsError } = useCollaborators();
  const updateProfile = useUpdateProfile();
  const updateOrganization = useUpdateOrganization();
  const updateCollaborator = useUpdateCollaborator();
  const [isEditing, setIsEditing] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [activeTab, setActiveTab] = useState("personal");
  
  // Get the user's organization or collaborator data
  const organization = organizations?.find(org => org.user_id === authUser?.id);
  const collaborator = collaborators?.find(collab => collab.user_id === authUser?.id);
  const isLoading = profileLoading || orgsLoading || collabsLoading;
  const hasError = profileError || orgsError || collabsError;

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-8">
        {/* Header Skeleton */}
        <div className="modern-card p-8 animate-pulse">
          <div className="flex items-center space-x-6">
            <div className="w-24 h-24 bg-slate-200 rounded-full"></div>
            <div className="space-y-3">
              <div className="h-6 bg-slate-200 rounded w-48"></div>
              <div className="h-4 bg-slate-200 rounded w-32"></div>
              <div className="h-4 bg-slate-200 rounded w-24"></div>
            </div>
          </div>
        </div>

        {/* Tabs Skeleton */}
        <div className="flex space-x-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="h-10 bg-slate-200 rounded w-24"></div>
          ))}
        </div>

        {/* Content Skeleton */}
        <div className="modern-card p-6 animate-pulse">
          <div className="space-y-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="space-y-2">
                <div className="h-4 bg-slate-200 rounded w-24"></div>
                <div className="h-10 bg-slate-200 rounded w-full"></div>
              </div>
            ))}
          </div>
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
          <h2 className="text-2xl font-bold text-slate-900 mb-2">Unable to Load Profile</h2>
          <p className="text-slate-600 mb-6">
            {profileError?.message || orgsError?.message || collabsError?.message || 
             "There was an error loading your profile. Please try again later."}
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

  const [formData, setFormData] = useState({
    name: "",
    phone: "",
    pincode: "",
    city: "",
    district: "",
    state: "",
    address: "",
    isRemoteOrg: false,
    orgType: "",
    purposeFocus: [],
    lastMonthPeopleFed: 0,
    collabType: ""
  });

  useEffect(() => {
    // Initialize form data with profile data when it loads
    if (profile) {
      setFormData({
        name: (profile.first_name || '') + ' ' + (profile.last_name || ''),
        phone: profile.phone || '',
        pincode: profile.pincode || '',
        city: profile.city || '',
        district: profile.district || '',
        state: profile.state || '',
        address: profile.address || '',
        isRemoteOrg: profile.is_remote_org || false,
        orgType: organization?.org_type || '',
        purposeFocus: organization?.purpose_focus || [],
        lastMonthPeopleFed: organization?.last_month_people_fed || 0,
        collabType: collaborator?.collab_type || ''
      });
    }
  }, [profile, organization, collaborator]);

  const handleEdit = () => {
    setFormData({
      name: (profile?.first_name || '') + ' ' + (profile?.last_name || ''),
      phone: profile?.phone || "",
      pincode: profile?.pincode || "",
      city: profile?.city || "",
      district: profile?.district || "",
      state: profile?.state || "",
      address: profile?.address || "",
      isRemoteOrg: profile?.is_remote_org || false,
      orgType: organization?.org_type || "",
      purposeFocus: organization?.purpose_focus || [],
      lastMonthPeopleFed: organization?.last_month_people_fed || 0,
      collabType: collaborator?.collab_type || ""
    });
    setIsEditing(true);
  };

  const handleSave = async () => {
    setIsSaving(true);
    setMessage("");
    
    try {
      // Split name into first and last name
      const [firstName, ...lastNameParts] = formData.name.split(' ');
      const lastName = lastNameParts.join(' ');
      
      // Update profile
      const profileData = {
        first_name: firstName,
        last_name: lastName,
        phone: formData.phone,
        pincode: formData.pincode,
        city: formData.city,
        district: formData.district,
        state: formData.state,
        address: formData.address,
        is_remote_org: formData.isRemoteOrg
      };
      
      await updateProfile.mutateAsync(profileData);
      
      // Update organization if user is ORG
      if (authUser?.role === "ORG" && organization) {
        const orgData = {
          org_type: formData.orgType,
          purpose_focus: formData.purposeFocus,
          last_month_people_fed: formData.lastMonthPeopleFed
        };
        await updateOrganization.mutateAsync({ id: organization.id, data: orgData });
      }
      
      // Update collaborator if user is COLLAB
      if (authUser?.role === "COLLAB" && collaborator) {
        const collabData = {
          collab_type: formData.collabType
        };
        await updateCollaborator.mutateAsync({ id: collaborator.id, data: collabData });
      }
      
      setIsEditing(false);
      setMessage("Profile updated successfully!");
      setTimeout(() => setMessage(""), 3000);
    } catch (error) {
      console.error('Profile update error:', error);
      setMessage(error.message || "Failed to update profile. Please try again.");
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
    setMessage("");
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === "checkbox" ? checked : value,
    });
  };

  const handlePurposeFocusChange = (purpose) => {
    setFormData({
      ...formData,
      purposeFocus: formData.purposeFocus.includes(purpose)
        ? formData.purposeFocus.filter(p => p !== purpose)
        : [...formData.purposeFocus, purpose]
    });
  };

  const purposeOptions = [
    { value: "CHILDREN", label: "Children" },
    { value: "ELDERLY", label: "Elderly" },
    { value: "WOMEN", label: "Women" },
    { value: "GENERAL", label: "General" },
    { value: "EMERGENCY", label: "Emergency" }
  ];

  if (isLoading) {
    return (
      <div className="space-y-8">
        {/* Profile Header Skeleton */}
        <div className="modern-card p-8 animate-pulse">
          <div className="flex items-start space-x-6">
            <div className="w-32 h-32 bg-slate-200 rounded-2xl"></div>
            <div className="flex-1">
              <div className="h-8 bg-slate-200 rounded w-1/3 mb-3"></div>
              <div className="h-4 bg-slate-200 rounded w-1/2 mb-4"></div>
              <div className="grid grid-cols-4 gap-4">
                <div className="h-16 bg-slate-200 rounded-xl"></div>
                <div className="h-16 bg-slate-200 rounded-xl"></div>
                <div className="h-16 bg-slate-200 rounded-xl"></div>
                <div className="h-16 bg-slate-200 rounded-xl"></div>
              </div>
            </div>
          </div>
        </div>
        
        {/* Content Skeleton */}
        <div className="modern-card p-6 animate-pulse">
          <div className="h-6 bg-slate-200 rounded w-1/4 mb-6"></div>
          <div className="grid grid-cols-2 gap-6">
            <div className="h-20 bg-slate-200 rounded-xl"></div>
            <div className="h-20 bg-slate-200 rounded-xl"></div>
            <div className="h-20 bg-slate-200 rounded-xl"></div>
            <div className="h-20 bg-slate-200 rounded-xl"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6 lg:space-y-8">
      {/* Profile Header */}
      <div className="modern-card p-6 lg:p-8 bg-gradient-to-br from-purple-50 to-pink-50 border-purple-200 animate-fade-in">
        <div className="flex flex-col lg:flex-row lg:items-start gap-6 lg:gap-8">
          {/* Profile Picture & Basic Info */}
          <div className="flex flex-col items-center lg:items-start flex-shrink-0">
            <div className="relative group">
              <div className="w-24 h-24 lg:w-32 lg:h-32 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl lg:rounded-2xl flex items-center justify-center text-white text-2xl lg:text-4xl font-bold shadow-2xl">
                {(profile?.first_name || authUser?.email)?.charAt(0)}
              </div>
              {isEditing && (
                <button className="absolute inset-0 bg-black/50 rounded-xl lg:rounded-2xl flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                  <Camera className="w-6 h-6 lg:w-8 lg:h-8 text-white" />
                </button>
              )}
              {profile?.verified && (
                <div className="absolute -top-2 -right-2 w-8 h-8 bg-emerald-500 rounded-full flex items-center justify-center">
                  <CheckCircle className="w-5 h-5 text-white" />
                </div>
              )}
            </div>
            <div className="text-center lg:text-left mt-4">
              <h1 className="text-3xl font-bold text-slate-900">{profile?.name}</h1>
              <p className="text-purple-600 font-medium text-lg">
                {authUser?.role === "ORG" ? "Organization" : authUser?.role === "COLLAB" ? "Food Donor" : "Administrator"}
              </p>
              {profile?.verified && (
                <div className="flex items-center justify-center lg:justify-start space-x-2 mt-2">
                  <Shield className="w-4 h-4 text-emerald-600" />
                  <span className="text-emerald-600 font-medium text-sm">Verified Account</span>
                </div>
              )}
            </div>
          </div>

          {/* Stats Grid */}
          <div className="flex-1">
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
              <div className="bg-white/80 p-4 rounded-xl text-center">
                <div className="text-2xl font-bold text-slate-900">{profile?.socialImpact?.mealsProvided || 0}</div>
                <div className="text-sm text-slate-600">Meals Provided</div>
              </div>
              <div className="bg-white/80 p-4 rounded-xl text-center">
                <div className="text-2xl font-bold text-emerald-600">{profile?.socialImpact?.foodWasteSaved || 0} kg</div>
                <div className="text-sm text-slate-600">Food Saved</div>
              </div>
              <div className="bg-white/80 p-4 rounded-xl text-center">
                <div className="text-2xl font-bold text-purple-600">{profile?.socialImpact?.communitiesServed || 0}</div>
                <div className="text-sm text-slate-600">Communities</div>
              </div>
              <div className="bg-white/80 p-4 rounded-xl text-center">
                <div className="text-2xl font-bold text-amber-600">{profile?.completionPercentage || 0}%</div>
                <div className="text-sm text-slate-600">Complete</div>
              </div>
            </div>

            {/* Quick Info */}
            <div className="bg-white/80 p-4 rounded-xl mb-6">
              <div className="flex items-center space-x-4 text-sm text-slate-600">
                <div className="flex items-center space-x-1">
                  <Calendar className="w-4 h-4" />
                  <span>Joined {new Date(profile?.joinedDate).toLocaleDateString()}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <MapPin className="w-4 h-4" />
                  <span>{profile?.city}, {profile?.state}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Mail className="w-4 h-4" />
                  <span>{authUser?.email}</span>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-wrap gap-3">
              {!isEditing ? (
                <button
                  onClick={handleEdit}
                  className="btn-modern px-6 py-3 inline-flex items-center group"
                >
                  <Edit className="h-5 w-5 mr-2 group-hover:scale-110 transition-transform" />
                  Edit Profile
                </button>
              ) : (
                <div className="flex space-x-3">
                  <button
                    onClick={handleSave}
                    disabled={isSaving}
                    className="btn-modern bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-600 hover:to-emerald-700 px-6 py-3 inline-flex items-center disabled:opacity-50"
                  >
                    <Save className="h-5 w-5 mr-2" />
                    {isSaving ? "Saving..." : "Save Changes"}
                  </button>
                  <button
                    onClick={handleCancel}
                    className="px-6 py-3 border-2 border-slate-200 text-slate-600 hover:border-emerald-300 hover:text-emerald-600 hover:bg-emerald-50 rounded-xl font-semibold transition-all"
                  >
                    <X className="h-5 w-5 mr-2 inline" />
                    Cancel
                  </button>
                </div>
              )}
              
              <button className="px-6 py-3 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-xl font-semibold transition-all">
                <Eye className="h-5 w-5 mr-2 inline" />
                View Public Profile
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Success Message */}
      {message && (
        <div className="modern-card bg-gradient-to-r from-emerald-50 to-emerald-100 border-emerald-200 p-4 animate-slide-in">
          <div className="flex items-center space-x-3">
            <div className="w-8 h-8 bg-emerald-500 rounded-full flex items-center justify-center">
              <CheckCircle className="h-5 w-5 text-white" />
            </div>
            <span className="text-emerald-700 font-medium">{message}</span>
          </div>
        </div>
      )}

      {/* Tab Navigation */}
      <div className="modern-card p-2 animate-slide-in">
        <nav className="flex space-x-2">
          {[
            { id: "personal", label: "Personal Info", icon: User },
            { id: "organization", label: authUser?.role === "ORG" ? "Organization" : "Business", icon: Building, show: authUser?.role !== "ADMIN" },
            { id: "achievements", label: "Achievements", icon: Award },
            { id: "settings", label: "Settings", icon: Shield }
          ].filter(tab => tab.show !== false).map((tab) => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center space-x-2 px-4 py-3 rounded-xl font-semibold text-sm transition-all ${
                  activeTab === tab.id
                    ? "bg-gradient-primary text-white shadow-lg"
                    : "text-slate-600 hover:bg-slate-100 hover:text-emerald-600"
                }`}
              >
                <Icon className="w-4 h-4" />
                <span>{tab.label}</span>
              </button>
            );
          })}
        </nav>
      </div>

      {/* Tab Content */}
      {activeTab === "personal" && (
        <div className="modern-card p-8 animate-fade-in">
          <h2 className="text-2xl font-bold text-slate-900 mb-6">Personal Information</h2>
          
          {isEditing ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">Full Name</label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your full name"
                />
              </div>
              
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">Phone Number</label>
                <input
                  type="tel"
                  name="phone"
                  value={formData.phone}
                  onChange={handleChange}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your phone number"
                />
              </div>
              
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">Pincode</label>
                <input
                  type="text"
                  name="pincode"
                  value={formData.pincode}
                  onChange={handleChange}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your pincode"
                />
              </div>
              
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">City</label>
                <input
                  type="text"
                  name="city"
                  value={formData.city}
                  onChange={handleChange}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your city"
                />
              </div>
              
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">District</label>
                <input
                  type="text"
                  name="district"
                  value={formData.district}
                  onChange={handleChange}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your district"
                />
              </div>
              
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">State</label>
                <input
                  type="text"
                  name="state"
                  value={formData.state}
                  onChange={handleChange}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your state"
                />
              </div>
              
              <div className="md:col-span-2">
                <label className="block text-sm font-semibold text-slate-700 mb-2">Address</label>
                <textarea
                  name="address"
                  value={formData.address}
                  onChange={handleChange}
                  rows={4}
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                  placeholder="Enter your complete address"
                />
              </div>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div className="space-y-6">
                <div className="flex items-center space-x-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center">
                    <User className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-500">Full Name</p>
                    <p className="text-lg font-semibold text-slate-900">{profile?.name || "Not provided"}</p>
                  </div>
                </div>
                
                <div className="flex items-center space-x-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center">
                    <Phone className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-500">Phone Number</p>
                    <p className="text-lg font-semibold text-slate-900">{profile?.phone || "Not provided"}</p>
                  </div>
                </div>
                
                <div className="flex items-center space-x-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-purple-400 to-purple-600 rounded-xl flex items-center justify-center">
                    <Mail className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-500">Email Address</p>
                    <p className="text-lg font-semibold text-slate-900">{authUser?.email}</p>
                  </div>
                </div>
              </div>
              
              <div className="space-y-6">
                <div className="flex items-center space-x-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-amber-400 to-amber-600 rounded-xl flex items-center justify-center">
                    <MapPin className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-500">Location</p>
                    <p className="text-lg font-semibold text-slate-900">
                      {profile?.city && profile?.state 
                        ? `${profile.city}, ${profile.state} - ${profile.pincode}`
                        : "Not provided"
                      }
                    </p>
                  </div>
                </div>
                
                <div className="flex items-center space-x-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-red-400 to-red-600 rounded-xl flex items-center justify-center">
                    <Globe className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-500">Account Type</p>
                    <p className="text-lg font-semibold text-slate-900">{authUser?.role}</p>
                  </div>
                </div>
                
                <div className="bg-slate-50 p-4 rounded-xl">
                  <p className="text-sm font-medium text-slate-500 mb-2">Complete Address</p>
                  <p className="text-slate-900">{profile?.address || "Not provided"}</p>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Organization/Business Tab */}
      {activeTab === "organization" && authUser?.role !== "ADMIN" && (
        <div className="modern-card p-8 animate-fade-in">
          <h2 className="text-2xl font-bold text-slate-900 mb-6">
            {authUser?.role === "ORG" ? "Organization Details" : "Business Information"}
          </h2>
          
          {authUser?.role === "ORG" ? (
            <div className="space-y-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <div className="space-y-6">
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center">
                      <Building className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Organization Name</p>
                      <p className="text-lg font-semibold text-slate-900">{organization?.orgName}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center">
                      <Shield className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Organization Type</p>
                      <p className="text-lg font-semibold text-slate-900">{organization?.orgType}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-purple-400 to-purple-600 rounded-xl flex items-center justify-center">
                      <Calendar className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Established</p>
                      <p className="text-lg font-semibold text-slate-900">{organization?.yearEstablished}</p>
                    </div>
                  </div>
                </div>
                
                <div className="space-y-6">
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-amber-400 to-amber-600 rounded-xl flex items-center justify-center">
                      <Users className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Team Size</p>
                      <p className="text-lg font-semibold text-slate-900">{organization?.teamSize} members</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center">
                      <Target className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">People Fed (Last Month)</p>
                      <p className="text-lg font-semibold text-slate-900">{organization?.lastMonthPeopleFed}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-red-400 to-red-600 rounded-xl flex items-center justify-center">
                      <Globe className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Website</p>
                      <a href={organization?.website} className="text-lg font-semibold text-blue-600 hover:text-blue-500 transition-colors">
                        {organization?.website}
                      </a>
                    </div>
                  </div>
                </div>
              </div>
              
              <div className="bg-slate-50 p-6 rounded-xl">
                <h3 className="text-lg font-bold text-slate-900 mb-3">Purpose Focus</h3>
                <div className="flex flex-wrap gap-3">
                  {organization?.purposeFocus?.map((purpose) => (
                    <span key={purpose} className="px-4 py-2 bg-gradient-to-r from-emerald-400 to-emerald-600 text-white font-medium rounded-xl">
                      {purposeOptions.find(opt => opt.value === purpose)?.label || purpose}
                    </span>
                  ))}
                </div>
              </div>
              
              <div className="bg-slate-50 p-6 rounded-xl">
                <h3 className="text-lg font-bold text-slate-900 mb-3">About Organization</h3>
                <p className="text-slate-600 leading-relaxed">{organization?.description}</p>
              </div>
            </div>
          ) : (
            <div className="space-y-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <div className="space-y-6">
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center">
                      <Building className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Business Name</p>
                      <p className="text-lg font-semibold text-slate-900">{collaborator?.businessName}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center">
                      <Target className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Business Type</p>
                      <p className="text-lg font-semibold text-slate-900">{collaborator?.collabType}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-amber-400 to-amber-600 rounded-xl flex items-center justify-center">
                      <Star className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Rating</p>
                      <p className="text-lg font-semibold text-slate-900">{collaborator?.rating} ⭐</p>
                    </div>
                  </div>
                </div>
                
                <div className="space-y-6">
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-purple-400 to-purple-600 rounded-xl flex items-center justify-center">
                      <Activity className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Monthly Donations (Avg)</p>
                      <p className="text-lg font-semibold text-slate-900">{collaborator?.averageMonthlyDonations}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center space-x-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-xl flex items-center justify-center">
                      <Zap className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <p className="text-sm font-medium text-slate-500">Tokens Earned</p>
                      <p className="text-lg font-semibold text-slate-900">{collaborator?.tokensEarned}</p>
                    </div>
                  </div>
                </div>
              </div>
              
              <div className="bg-slate-50 p-6 rounded-xl">
                <h3 className="text-lg font-bold text-slate-900 mb-3">Specializations</h3>
                <div className="flex flex-wrap gap-3">
                  {collaborator?.specializations?.map((spec) => (
                    <span key={spec} className="px-4 py-2 bg-gradient-to-r from-blue-400 to-blue-600 text-white font-medium rounded-xl">
                      {spec}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Achievements Tab */}
      {activeTab === "achievements" && (
        <div className="modern-card p-8 animate-fade-in">
          <h2 className="text-2xl font-bold text-slate-900 mb-6">Achievements & Recognition</h2>
          
          {/* Dynamic achievements based on user activity */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {/* Profile completeness achievement */}
            {profile?.first_name && profile?.phone && profile?.address && (
              <div className="bg-gradient-to-br from-emerald-50 to-green-50 p-6 rounded-xl border border-emerald-200">
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-emerald-400 to-green-500 rounded-xl flex items-center justify-center">
                    <CheckCircle className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-bold text-slate-900">Profile Complete</h3>
                    <p className="text-emerald-600 font-medium">Your profile is fully set up</p>
                  </div>
                </div>
              </div>
            )}
            
            {/* Organization membership achievement */}
            {authUser?.role === "ORG" && organization && (
              <div className="bg-gradient-to-br from-blue-50 to-cyan-50 p-6 rounded-xl border border-blue-200">
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-cyan-500 rounded-xl flex items-center justify-center">
                    <Building className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-bold text-slate-900">Organization Member</h3>
                    <p className="text-blue-600 font-medium">Registered organization</p>
                  </div>
                </div>
              </div>
            )}
            
            {/* Collaborator achievement */}
            {authUser?.role === "COLLAB" && collaborator && (
              <div className="bg-gradient-to-br from-purple-50 to-pink-50 p-6 rounded-xl border border-purple-200">
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-purple-400 to-pink-500 rounded-xl flex items-center justify-center">
                    <Users className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-bold text-slate-900">Community Collaborator</h3>
                    <p className="text-purple-600 font-medium">Contributing to food security</p>
                  </div>
                </div>
              </div>
            )}
            
            {/* Platform newcomer fallback */}
            {!profile?.first_name && (
              <div className="bg-gradient-to-br from-amber-50 to-orange-50 p-6 rounded-xl border border-amber-200">
                <div className="flex items-center space-x-3 mb-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-amber-400 to-orange-500 rounded-xl flex items-center justify-center">
                    <Sparkles className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h3 className="font-bold text-slate-900">Welcome to FoodFlow!</h3>
                    <p className="text-amber-600 font-medium">Complete your profile to earn more achievements</p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Settings Tab */}
      {activeTab === "settings" && (
        <div className="modern-card p-8 animate-fade-in">
          <h2 className="text-2xl font-bold text-slate-900 mb-6">Account Settings</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="space-y-6">
              <h3 className="text-lg font-bold text-slate-900">Security</h3>
              
              <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl">
                <div className="flex items-center space-x-3">
                  <Shield className="w-5 h-5 text-emerald-600" />
                  <div>
                    <p className="font-medium text-slate-900">Account Verification</p>
                    <p className="text-sm text-slate-600">Your account is verified</p>
                  </div>
                </div>
                <div className="w-8 h-8 bg-emerald-500 rounded-full flex items-center justify-center">
                  <CheckCircle className="w-5 h-5 text-white" />
                </div>
              </div>
              
              <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl">
                <div className="flex items-center space-x-3">
                  <Mail className="w-5 h-5 text-blue-600" />
                  <div>
                    <p className="font-medium text-slate-900">Email Notifications</p>
                    <p className="text-sm text-slate-600">Receive updates about offers</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" className="sr-only peer" defaultChecked />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                </label>
              </div>
            </div>
            
            <div className="space-y-6">
              <h3 className="text-lg font-bold text-slate-900">Privacy</h3>
              
              <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl">
                <div className="flex items-center space-x-3">
                  <Eye className="w-5 h-5 text-purple-600" />
                  <div>
                    <p className="font-medium text-slate-900">Public Profile</p>
                    <p className="text-sm text-slate-600">Show profile to other users</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" className="sr-only peer" defaultChecked />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-purple-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
                </label>
              </div>
              
              <div className="flex items-center justify-between p-4 bg-slate-50 rounded-xl">
                <div className="flex items-center space-x-3">
                  <Activity className="w-5 h-5 text-emerald-600" />
                  <div>
                    <p className="font-medium text-slate-900">Activity Tracking</p>
                    <p className="text-sm text-slate-600">Track donation activities</p>
                  </div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" className="sr-only peer" defaultChecked />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-emerald-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-emerald-600"></div>
                </label>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
