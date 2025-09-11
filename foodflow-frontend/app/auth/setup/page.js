"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, Building2, User, Heart, Users, Baby, AlertTriangle, MapPin, Phone, Mail } from "lucide-react";
import { useAuth } from "../../../lib/auth";
import { organizationsApi, collaboratorsApi } from "../../../lib/api";

export default function ProfileSetupPage() {
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const router = useRouter();
  
  const [step, setStep] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [formData, setFormData] = useState({
    // Organization fields
    orgType: "",
    purposeFocus: [],
    lastMonthPeopleFed: "",
    // Collaborator fields
    collabType: "",
    // Profile fields
    firstName: "",
    lastName: "",
    phone: "",
    pincode: "",
    city: "",
    district: "",
    state: "",
    address: "",
    isRemoteOrg: false
  });

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/auth/login");
    }
  }, [isAuthenticated, authLoading, router]);

  if (authLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-emerald-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-emerald-600"></div>
      </div>
    );
  }

  if (!isAuthenticated || !user) {
    return null;
  }

  const purposeOptions = [
    { value: "CHILDREN", label: "Children", icon: Baby, desc: "Programs focused on children" },
    { value: "ELDERLY", label: "Elderly", icon: Heart, desc: "Support for elderly community" },
    { value: "WOMEN", label: "Women", icon: Users, desc: "Women empowerment programs" },
    { value: "GENERAL", label: "General", icon: Users, desc: "General community support" },
    { value: "EMERGENCY", label: "Emergency", icon: AlertTriangle, desc: "Emergency relief efforts" }
  ];

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value
    }));
  };

  const handlePurposeFocusChange = (value) => {
    setFormData(prev => ({
      ...prev,
      purposeFocus: prev.purposeFocus.includes(value)
        ? prev.purposeFocus.filter(p => p !== value)
        : [...prev.purposeFocus, value]
    }));
  };

  const validateStep1 = () => {
    if (user.role === "ORG") {
      return formData.orgType && formData.purposeFocus.length > 0 && formData.lastMonthPeopleFed;
    } else if (user.role === "COLLAB") {
      return formData.collabType;
    }
    return false;
  };

  const validateStep2 = () => {
    const isValidPhone = formData.phone && formData.phone.length === 10 && /^\d{10}$/.test(formData.phone);
    const isValidPincode = formData.pincode && formData.pincode.length === 6 && /^\d{6}$/.test(formData.pincode);
    const isValidName = formData.firstName && formData.lastName && 
                       formData.firstName.length >= 2 && formData.lastName.length >= 2;
    const isValidLocation = formData.city && formData.state && 
                           formData.city.length >= 2 && formData.state.length >= 2;
    
    return isValidName && isValidPhone && isValidPincode && isValidLocation;
  };

  const handleNext = () => {
    if (step === 1 && validateStep1()) {
      setStep(2);
    }
  };

  const handleBack = () => {
    if (step === 2) {
      setStep(1);
    }
  };

  const handleSubmit = async () => {
    if (!validateStep2()) {
      let errorMsg = "Please fix the following issues:\n";
      if (!formData.firstName || formData.firstName.length < 2) errorMsg += "• First name must be at least 2 characters\n";
      if (!formData.lastName || formData.lastName.length < 2) errorMsg += "• Last name must be at least 2 characters\n";
      if (!formData.phone || formData.phone.length !== 10 || !/^\d{10}$/.test(formData.phone)) errorMsg += "• Phone must be exactly 10 digits\n";
      if (!formData.pincode || formData.pincode.length !== 6 || !/^\d{6}$/.test(formData.pincode)) errorMsg += "• PIN code must be exactly 6 digits\n";
      if (!formData.city || formData.city.length < 2) errorMsg += "• City must be at least 2 characters\n";
      if (!formData.state || formData.state.length < 2) errorMsg += "• State must be at least 2 characters\n";
      
      setError(errorMsg);
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      console.log("Starting profile setup for user:", user);
      
      // Step 1: Update profile
      const profileData = {
        name: `${formData.firstName} ${formData.lastName}`.trim(),
        phone: formData.phone,
        pincode: formData.pincode,
        city: formData.city,
        district: formData.district || "",
        state: formData.state,
        address: formData.address || "",
        is_remote_org: formData.isRemoteOrg
      };

      console.log("Updating profile with data:", profileData);

      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/v1/profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`
        },
        body: JSON.stringify(profileData)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        console.error("Profile update failed:", response.status, errorData);
        throw new Error(errorData.error?.message || 'Failed to update profile');
      }

      const profileResult = await response.json();
      console.log("Profile updated successfully:", profileResult);

      // Step 2: Create organization or collaborator profile
      if (user.role === "ORG") {
        const orgData = {
          org_type: formData.orgType,
          purpose_focus: formData.purposeFocus,
          last_month_people_fed: parseInt(formData.lastMonthPeopleFed)
        };
        console.log("Creating organization with data:", orgData);
        
        try {
          const orgResult = await organizationsApi.create(orgData);
          console.log("Organization created successfully:", orgResult);
        } catch (orgError) {
          console.error("Organization creation failed:", orgError);
          throw orgError;
        }
      } else if (user.role === "COLLAB") {
        const collabData = {
          collab_type: formData.collabType
        };
        console.log("Creating collaborator with data:", collabData);
        
        try {
          const collabResult = await collaboratorsApi.create(collabData);
          console.log("Collaborator created successfully:", collabResult);
        } catch (collabError) {
          console.error("Collaborator creation failed:", collabError);
          throw collabError;
        }
      }

      console.log("Profile setup completed successfully, redirecting to dashboard");
      // Redirect to dashboard
      router.push("/dashboard");
    } catch (err) {
      console.error("Profile setup error:", err);
      setError(err.message || "Failed to complete profile setup. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-emerald-50 relative overflow-hidden flex items-center justify-center py-12">
      {/* Background decorations */}
      <div className="absolute inset-0 bg-hero-pattern opacity-30"></div>
      <div className="absolute top-10 left-10 w-72 h-72 bg-gradient-to-r from-emerald-400 to-cyan-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-float"></div>
      <div className="absolute bottom-10 right-10 w-72 h-72 bg-gradient-to-r from-purple-400 to-pink-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-float" style={{animationDelay: '2s'}}></div>
      
      <div className="w-full max-w-2xl relative z-10">
        {/* Header */}
        <div className="text-center mb-8 animate-fade-in">
          <h2 className="text-4xl font-bold text-slate-900 mb-3">
            Complete Your Profile
          </h2>
          <p className="text-slate-600 text-lg">
            Let's set up your {user.role === "ORG" ? "organization" : "collaborator"} profile
          </p>
          <div className="flex justify-center mt-6">
            <div className="flex items-center space-x-4">
              <div className={`w-10 h-10 rounded-full flex items-center justify-center ${step >= 1 ? 'bg-emerald-500 text-white' : 'bg-slate-200 text-slate-400'}`}>
                {user.role === "ORG" ? <Building2 className="w-5 h-5" /> : <User className="w-5 h-5" />}
              </div>
              <div className={`w-16 h-1 ${step >= 2 ? 'bg-emerald-500' : 'bg-slate-200'}`}></div>
              <div className={`w-10 h-10 rounded-full flex items-center justify-center ${step >= 2 ? 'bg-emerald-500 text-white' : 'bg-slate-200 text-slate-400'}`}>
                <MapPin className="w-5 h-5" />
              </div>
            </div>
          </div>
        </div>

        {/* Form */}
        <div className="modern-card p-8 animate-slide-in">
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-xl text-sm mb-6">
              <pre className="whitespace-pre-wrap font-sans">{error}</pre>
            </div>
          )}

          {step === 1 && (
            <div className="space-y-6">
              <h3 className="text-2xl font-bold text-slate-900 mb-6">
                {user.role === "ORG" ? "Organization Details" : "Collaborator Details"}
              </h3>

              {user.role === "ORG" ? (
                <>
                  {/* Organization Type */}
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-2">
                      Organization Type *
                    </label>
                    <input
                      type="text"
                      name="orgType"
                      value={formData.orgType}
                      onChange={handleInputChange}
                      className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                      placeholder="e.g., NGO, Charity, Community Kitchen, School"
                    />
                  </div>

                  {/* Purpose Focus */}
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-3">
                      Purpose Focus * (Select all that apply)
                    </label>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                      {purposeOptions.map((option) => {
                        const Icon = option.icon;
                        return (
                          <button
                            key={option.value}
                            type="button"
                            onClick={() => handlePurposeFocusChange(option.value)}
                            className={`p-4 rounded-xl border-2 text-left transition-all ${
                              formData.purposeFocus.includes(option.value)
                                ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                                : 'border-slate-200 hover:border-emerald-300 hover:bg-emerald-50/50'
                            }`}
                          >
                            <div className="flex items-start space-x-3">
                              <Icon className="w-5 h-5 mt-0.5 flex-shrink-0" />
                              <div>
                                <div className="font-semibold">{option.label}</div>
                                <div className="text-sm opacity-75">{option.desc}</div>
                              </div>
                            </div>
                          </button>
                        );
                      })}
                    </div>
                  </div>

                  {/* People Fed Last Month */}
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-2">
                      People Fed Last Month *
                    </label>
                    <input
                      type="number"
                      name="lastMonthPeopleFed"
                      value={formData.lastMonthPeopleFed}
                      onChange={handleInputChange}
                      min="0"
                      className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                      placeholder="Number of people your organization fed last month"
                    />
                  </div>
                </>
              ) : (
                <>
                  {/* Collaborator Type */}
                  <div>
                    <label className="block text-sm font-semibold text-slate-700 mb-2">
                      Collaborator Type *
                    </label>
                    <select
                      name="collabType"
                      value={formData.collabType}
                      onChange={handleInputChange}
                      className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    >
                      <option value="">Select collaborator type</option>
                      <option value="Restaurant">Restaurant</option>
                      <option value="Cafe">Cafe</option>
                      <option value="Hotel">Hotel</option>
                      <option value="Catering Service">Catering Service</option>
                      <option value="Food Store">Food Store</option>
                      <option value="Individual Donor">Individual Donor</option>
                      <option value="Food Production">Food Production</option>
                      <option value="Other">Other</option>
                    </select>
                  </div>
                </>
              )}

              <div className="flex justify-end pt-4">
                <button
                  onClick={handleNext}
                  disabled={!validateStep1()}
                  className="btn-modern px-8 py-3 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next Step
                </button>
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <h3 className="text-2xl font-bold text-slate-900">
                  Contact & Location Details
                </h3>
                <button
                  onClick={handleBack}
                  className="flex items-center space-x-2 text-slate-600 hover:text-emerald-600 transition-colors"
                >
                  <ArrowLeft className="w-4 h-4" />
                  <span>Back</span>
                </button>
              </div>

              {/* Name Fields */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">
                    First Name *
                  </label>
                  <input
                    type="text"
                    name="firstName"
                    value={formData.firstName}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="First name"
                  />
                </div>
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">
                    Last Name *
                  </label>
                  <input
                    type="text"
                    name="lastName"
                    value={formData.lastName}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="Last name"
                  />
                </div>
              </div>

              {/* Phone */}
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">
                  Phone Number *
                </label>
                <div className="relative">
                  <Phone className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-slate-400" />
                  <input
                    type="tel"
                    name="phone"
                    value={formData.phone}
                    onChange={handleInputChange}
                    className="w-full pl-10 pr-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="10-digit phone number"
                    maxLength="10"
                  />
                </div>
                {formData.phone && formData.phone.length !== 10 && (
                  <p className="mt-1 text-sm text-red-600">Phone number must be exactly 10 digits</p>
                )}
              </div>

              {/* Location Fields */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">
                    PIN Code *
                  </label>
                  <input
                    type="text"
                    name="pincode"
                    value={formData.pincode}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="6-digit PIN code"
                    maxLength="6"
                  />
                  {formData.pincode && formData.pincode.length !== 6 && (
                    <p className="mt-1 text-sm text-red-600">PIN code must be exactly 6 digits</p>
                  )}
                </div>
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">
                    City *
                  </label>
                  <input
                    type="text"
                    name="city"
                    value={formData.city}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="City"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">
                    District
                  </label>
                  <input
                    type="text"
                    name="district"
                    value={formData.district}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="District"
                  />
                </div>
                <div>
                  <label className="block text-sm font-semibold text-slate-700 mb-2">
                    State *
                  </label>
                  <input
                    type="text"
                    name="state"
                    value={formData.state}
                    onChange={handleInputChange}
                    className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all"
                    placeholder="State"
                  />
                </div>
              </div>

              {/* Address */}
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-2">
                  Address
                </label>
                <textarea
                  name="address"
                  value={formData.address}
                  onChange={handleInputChange}
                  rows="3"
                  className="w-full px-4 py-3 border border-slate-200 rounded-xl placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 text-slate-900 bg-white/50 backdrop-blur-sm transition-all resize-none"
                  placeholder="Full address (optional)"
                />
              </div>

              {/* Remote Organization Checkbox (for ORG users) */}
              {user.role === "ORG" && (
                <div className="flex items-start">
                  <div className="flex items-center h-5">
                    <input
                      id="isRemoteOrg"
                      name="isRemoteOrg"
                      type="checkbox"
                      checked={formData.isRemoteOrg}
                      onChange={handleInputChange}
                      className="h-4 w-4 text-emerald-600 focus:ring-emerald-500 border-slate-300 rounded transition-colors"
                    />
                  </div>
                  <div className="ml-3 text-sm">
                    <label htmlFor="isRemoteOrg" className="font-medium text-slate-700">
                      We serve remote/underserved areas
                    </label>
                    <p className="text-slate-500">Check this if your organization primarily serves remote or underserved communities</p>
                  </div>
                </div>
              )}

              <div className="flex justify-end pt-4">
                <button
                  onClick={handleSubmit}
                  disabled={isLoading || !validateStep2()}
                  className="btn-modern px-8 py-3 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isLoading ? (
                    <div className="flex items-center">
                      <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></div>
                      Setting up...
                    </div>
                  ) : (
                    "Complete Setup"
                  )}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}