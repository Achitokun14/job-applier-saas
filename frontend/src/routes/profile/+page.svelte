<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Badge } from '$lib/components/ui/badge';
  import { User, Briefcase, GraduationCap, Code, FolderOpen, CheckCircle, AlertCircle, Loader2, Save, Plus, X, Globe, Linkedin, Phone, MapPin, Award, Languages, Trophy } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let loading = $state(true);
  let saving = $state(false);

  // Personal info
  let name = $state('');
  let email = $state('');
  let phone = $state('');
  let location = $state('');
  let website = $state('');
  let linkedin = $state('');
  let summary = $state('');

  // Experience
  let experience = $state([]);

  // Education
  let education = $state([]);

  // Skills
  let skills = $state([]);
  let skillInput = $state('');

  // Projects
  let projects = $state([]);

  // Certifications
  let certifications = $state([]);

  // Languages
  let languages = $state([]);

  // Achievements
  let achievements = $state('');

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    try {
      const token = $auth.token;
      const profile = await api.getProfile(token);

      name = profile.name || '';
      email = profile.email || '';

      // Parse personal_info JSON
      try {
        const personalInfo = typeof profile.personal_info === 'string'
          ? JSON.parse(profile.personal_info)
          : (profile.personal_info || {});
        phone = personalInfo.phone || '';
        location = personalInfo.location || '';
        website = personalInfo.website || '';
        linkedin = personalInfo.linkedin || '';
        summary = personalInfo.summary || '';
      } catch { /* ignore parse errors */ }

      // Parse experience JSON
      try {
        experience = typeof profile.experience === 'string'
          ? JSON.parse(profile.experience)
          : (profile.experience || []);
        if (!Array.isArray(experience)) experience = [];
      } catch { experience = []; }

      // Parse education JSON
      try {
        education = typeof profile.education === 'string'
          ? JSON.parse(profile.education)
          : (profile.education || []);
        if (!Array.isArray(education)) education = [];
      } catch { education = []; }

      // Parse skills (comma-separated string)
      try {
        if (typeof profile.skills === 'string' && profile.skills.trim()) {
          skills = profile.skills.split(',').map(s => s.trim()).filter(s => s);
        } else if (Array.isArray(profile.skills)) {
          skills = profile.skills;
        } else {
          skills = [];
        }
      } catch { skills = []; }

      // Parse projects JSON
      try {
        projects = typeof profile.projects === 'string'
          ? JSON.parse(profile.projects)
          : (profile.projects || []);
        if (!Array.isArray(projects)) projects = [];
      } catch { projects = []; }

      // Parse certifications JSON
      try {
        certifications = typeof profile.certifications === 'string'
          ? JSON.parse(profile.certifications)
          : (profile.certifications || []);
        if (!Array.isArray(certifications)) certifications = [];
      } catch { certifications = []; }

      // Parse languages JSON
      try {
        languages = typeof profile.languages === 'string'
          ? JSON.parse(profile.languages)
          : (profile.languages || []);
        if (!Array.isArray(languages)) languages = [];
      } catch { languages = []; }

      // Achievements plain text
      achievements = profile.achievements || '';

    } catch (e) {
      console.error(e);
      toast.error('Failed to load profile');
    } finally {
      loading = false;
    }
  });

  async function saveProfile() {
    saving = true;
    try {
      const token = $auth.token;
      await api.updateProfile({
        name,
        personal_info: JSON.stringify({ phone, location, website, linkedin, summary }),
        experience: JSON.stringify(experience),
        education: JSON.stringify(education),
        skills: skills.join(', '),
        projects: JSON.stringify(projects),
        certifications: JSON.stringify(certifications),
        languages: JSON.stringify(languages),
        achievements
      }, token);
      toast.success('Profile saved successfully!');
    } catch (e) {
      toast.error('Failed to save profile. Please try again.');
    } finally {
      saving = false;
    }
  }

  // Experience helpers
  function addExperience() {
    experience = [...experience, { title: '', company: '', startDate: '', endDate: '', current: false, description: '' }];
  }
  function removeExperience(index) {
    experience = experience.filter((_, i) => i !== index);
  }

  // Education helpers
  function addEducation() {
    education = [...education, { degree: '', school: '', startYear: '', endYear: '', field: '' }];
  }
  function removeEducation(index) {
    education = education.filter((_, i) => i !== index);
  }

  // Skill helpers
  function addSkill() {
    const trimmed = skillInput.trim();
    if (trimmed && !skills.includes(trimmed)) {
      skills = [...skills, trimmed];
      skillInput = '';
    }
  }
  function removeSkill(index) {
    skills = skills.filter((_, i) => i !== index);
  }
  function handleSkillKeydown(e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      addSkill();
    }
  }

  // Project helpers
  function addProject() {
    projects = [...projects, { name: '', description: '', technologies: '', url: '' }];
  }
  function removeProject(index) {
    projects = projects.filter((_, i) => i !== index);
  }

  // Certification helpers
  function addCertification() {
    certifications = [...certifications, { name: '', issuer: '', date: '', url: '' }];
  }
  function removeCertification(index) {
    certifications = certifications.filter((_, i) => i !== index);
  }

  // Language helpers
  function addLanguage() {
    languages = [...languages, { language: '', proficiency: 'intermediate' }];
  }
  function removeLanguage(index) {
    languages = languages.filter((_, i) => i !== index);
  }

  const textareaClass = "flex w-full rounded-lg border border-input bg-background px-3 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y transition-colors";
</script>

<svelte:head>
  <title>My Profile - JobApplier</title>
</svelte:head>

<div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Header -->
  <div class="mb-8 animate-fade-in">
    <h1 class="text-2xl sm:text-3xl font-bold text-foreground">My Profile</h1>
    <p class="text-muted-foreground mt-1">Build your structured CV for better AI-generated applications</p>
  </div>

  {#if loading}
    <div class="space-y-6 animate-fade-in">
      {#each [1,2,3,4,5,6] as _}
        <div class="skeleton h-48 rounded-xl"></div>
      {/each}
    </div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); saveProfile(); }} class="space-y-6">

      <!-- Personal Information -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center">
              <User size={16} class="text-primary" />
            </div>
            <div>
              <CardTitle class="text-base">Personal Information</CardTitle>
              <CardDescription class="text-xs">Your contact details and professional summary</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-2">
              <label for="name" class="text-sm font-medium text-foreground">Full Name</label>
              <Input type="text" id="name" bind:value={name} placeholder="John Doe" />
            </div>
            <div class="space-y-2">
              <label for="email" class="text-sm font-medium text-foreground">Email</label>
              <Input type="email" id="email" bind:value={email} placeholder="john@example.com" disabled />
            </div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-2">
              <label for="phone" class="text-sm font-medium text-foreground flex items-center gap-1.5">
                <Phone size={13} class="text-muted-foreground" />
                Phone
              </label>
              <Input type="tel" id="phone" bind:value={phone} placeholder="+1 (555) 123-4567" />
            </div>
            <div class="space-y-2">
              <label for="location" class="text-sm font-medium text-foreground flex items-center gap-1.5">
                <MapPin size={13} class="text-muted-foreground" />
                Location
              </label>
              <Input type="text" id="location" bind:value={location} placeholder="San Francisco, CA" />
            </div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-2">
              <label for="website" class="text-sm font-medium text-foreground flex items-center gap-1.5">
                <Globe size={13} class="text-muted-foreground" />
                Website
              </label>
              <Input type="url" id="website" bind:value={website} placeholder="https://example.com" />
            </div>
            <div class="space-y-2">
              <label for="linkedin" class="text-sm font-medium text-foreground flex items-center gap-1.5">
                <Linkedin size={13} class="text-muted-foreground" />
                LinkedIn
              </label>
              <Input type="url" id="linkedin" bind:value={linkedin} placeholder="https://linkedin.com/in/..." />
            </div>
          </div>
          <div class="space-y-2">
            <label for="summary" class="text-sm font-medium text-foreground">Professional Summary</label>
            <textarea
              id="summary"
              bind:value={summary}
              rows="3"
              placeholder="A brief professional summary highlighting your career goals and key strengths..."
              class={textareaClass}
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Experience -->
      <Card class="animate-fade-in-up delay-100">
        <CardHeader class="pb-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-lg bg-blue-500/10 flex items-center justify-center">
                <Briefcase size={16} class="text-blue-500" />
              </div>
              <div>
                <CardTitle class="text-base">Work Experience</CardTitle>
                <CardDescription class="text-xs">Your professional work history</CardDescription>
              </div>
            </div>
            <Button type="button" variant="outline" size="sm" onclick={addExperience}>
              <Plus size={14} class="mr-1.5" />
              Add
            </Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          {#if experience.length === 0}
            <p class="text-sm text-muted-foreground text-center py-4">No experience added yet. Click "Add" to add your work history.</p>
          {/if}
          {#each experience as exp, i}
            <div class="border border-border rounded-lg p-4 space-y-3 relative">
              <button
                type="button"
                onclick={() => removeExperience(i)}
                class="absolute top-3 right-3 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
                aria-label="Remove experience"
              >
                <X size={16} />
              </button>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pr-8">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Job Title</label>
                  <Input type="text" bind:value={experience[i].title} placeholder="Software Engineer" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Company</label>
                  <Input type="text" bind:value={experience[i].company} placeholder="Acme Corp" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Start Date</label>
                  <Input type="text" bind:value={experience[i].startDate} placeholder="Jan 2022" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">End Date</label>
                  <Input type="text" bind:value={experience[i].endDate} placeholder="Present" disabled={experience[i].current} />
                </div>
                <div class="flex items-end pb-1">
                  <label class="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" bind:checked={experience[i].current} class="h-4 w-4 rounded border-input text-primary focus:ring-ring" />
                    <span class="text-sm text-foreground">Current</span>
                  </label>
                </div>
              </div>
              <div class="space-y-1.5">
                <label class="text-xs font-medium text-muted-foreground">Description</label>
                <textarea
                  bind:value={experience[i].description}
                  rows="3"
                  placeholder="Key responsibilities and achievements..."
                  class={textareaClass}
                ></textarea>
              </div>
            </div>
          {/each}
        </CardContent>
      </Card>

      <!-- Education -->
      <Card class="animate-fade-in-up delay-200">
        <CardHeader class="pb-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-lg bg-green-500/10 flex items-center justify-center">
                <GraduationCap size={16} class="text-green-500" />
              </div>
              <div>
                <CardTitle class="text-base">Education</CardTitle>
                <CardDescription class="text-xs">Your academic background</CardDescription>
              </div>
            </div>
            <Button type="button" variant="outline" size="sm" onclick={addEducation}>
              <Plus size={14} class="mr-1.5" />
              Add
            </Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          {#if education.length === 0}
            <p class="text-sm text-muted-foreground text-center py-4">No education added yet. Click "Add" to add your academic history.</p>
          {/if}
          {#each education as edu, i}
            <div class="border border-border rounded-lg p-4 space-y-3 relative">
              <button
                type="button"
                onclick={() => removeEducation(i)}
                class="absolute top-3 right-3 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
                aria-label="Remove education"
              >
                <X size={16} />
              </button>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pr-8">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Degree</label>
                  <Input type="text" bind:value={education[i].degree} placeholder="Bachelor of Science" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">School</label>
                  <Input type="text" bind:value={education[i].school} placeholder="MIT" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Field of Study</label>
                  <Input type="text" bind:value={education[i].field} placeholder="Computer Science" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Start Year</label>
                  <Input type="text" bind:value={education[i].startYear} placeholder="2018" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">End Year</label>
                  <Input type="text" bind:value={education[i].endYear} placeholder="2022" />
                </div>
              </div>
            </div>
          {/each}
        </CardContent>
      </Card>

      <!-- Skills -->
      <Card class="animate-fade-in-up delay-300">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-amber-500/10 flex items-center justify-center">
              <Code size={16} class="text-amber-500" />
            </div>
            <div>
              <CardTitle class="text-base">Skills</CardTitle>
              <CardDescription class="text-xs">Your technical and professional skills</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="flex gap-2">
            <Input
              type="text"
              bind:value={skillInput}
              placeholder="Type a skill and press Enter..."
              onkeydown={handleSkillKeydown}
              class="flex-1"
            />
            <Button type="button" variant="outline" onclick={addSkill} disabled={!skillInput.trim()}>
              <Plus size={14} class="mr-1.5" />
              Add
            </Button>
          </div>
          {#if skills.length > 0}
            <div class="flex flex-wrap gap-2">
              {#each skills as skill, i}
                <Badge variant="secondary" class="text-sm px-3 py-1.5 flex items-center gap-1.5">
                  {skill}
                  <button
                    type="button"
                    onclick={() => removeSkill(i)}
                    class="text-muted-foreground hover:text-destructive transition-colors cursor-pointer ml-0.5"
                    aria-label="Remove skill"
                  >
                    <X size={12} />
                  </button>
                </Badge>
              {/each}
            </div>
          {:else}
            <p class="text-sm text-muted-foreground text-center py-2">No skills added yet. Type a skill and press Enter.</p>
          {/if}
        </CardContent>
      </Card>

      <!-- Projects -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-lg bg-purple-500/10 flex items-center justify-center">
                <FolderOpen size={16} class="text-purple-500" />
              </div>
              <div>
                <CardTitle class="text-base">Projects</CardTitle>
                <CardDescription class="text-xs">Notable projects you have worked on</CardDescription>
              </div>
            </div>
            <Button type="button" variant="outline" size="sm" onclick={addProject}>
              <Plus size={14} class="mr-1.5" />
              Add
            </Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          {#if projects.length === 0}
            <p class="text-sm text-muted-foreground text-center py-4">No projects added yet. Click "Add" to showcase your work.</p>
          {/if}
          {#each projects as proj, i}
            <div class="border border-border rounded-lg p-4 space-y-3 relative">
              <button
                type="button"
                onclick={() => removeProject(i)}
                class="absolute top-3 right-3 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
                aria-label="Remove project"
              >
                <X size={16} />
              </button>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pr-8">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Project Name</label>
                  <Input type="text" bind:value={projects[i].name} placeholder="My Awesome Project" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">URL</label>
                  <Input type="url" bind:value={projects[i].url} placeholder="https://github.com/..." />
                </div>
              </div>
              <div class="space-y-1.5">
                <label class="text-xs font-medium text-muted-foreground">Technologies</label>
                <Input type="text" bind:value={projects[i].technologies} placeholder="React, Node.js, PostgreSQL" />
              </div>
              <div class="space-y-1.5">
                <label class="text-xs font-medium text-muted-foreground">Description</label>
                <textarea
                  bind:value={projects[i].description}
                  rows="2"
                  placeholder="Brief description of the project and your contributions..."
                  class={textareaClass}
                ></textarea>
              </div>
            </div>
          {/each}
        </CardContent>
      </Card>

      <!-- Certifications -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-lg bg-rose-500/10 flex items-center justify-center">
                <Award size={16} class="text-rose-500" />
              </div>
              <div>
                <CardTitle class="text-base">Certifications</CardTitle>
                <CardDescription class="text-xs">Professional certifications and licenses</CardDescription>
              </div>
            </div>
            <Button type="button" variant="outline" size="sm" onclick={addCertification}>
              <Plus size={14} class="mr-1.5" />
              Add
            </Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          {#if certifications.length === 0}
            <p class="text-sm text-muted-foreground text-center py-4">No certifications added yet. Click "Add" to list your credentials.</p>
          {/if}
          {#each certifications as cert, i}
            <div class="border border-border rounded-lg p-4 space-y-3 relative">
              <button
                type="button"
                onclick={() => removeCertification(i)}
                class="absolute top-3 right-3 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
                aria-label="Remove certification"
              >
                <X size={16} />
              </button>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pr-8">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Certification Name</label>
                  <Input type="text" bind:value={certifications[i].name} placeholder="AWS Solutions Architect" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Issuer</label>
                  <Input type="text" bind:value={certifications[i].issuer} placeholder="Amazon Web Services" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Date</label>
                  <Input type="text" bind:value={certifications[i].date} placeholder="Mar 2024" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">URL</label>
                  <Input type="url" bind:value={certifications[i].url} placeholder="https://verify.example.com/..." />
                </div>
              </div>
            </div>
          {/each}
        </CardContent>
      </Card>

      <!-- Languages -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-lg bg-cyan-500/10 flex items-center justify-center">
                <Languages size={16} class="text-cyan-500" />
              </div>
              <div>
                <CardTitle class="text-base">Languages</CardTitle>
                <CardDescription class="text-xs">Languages you speak and your proficiency level</CardDescription>
              </div>
            </div>
            <Button type="button" variant="outline" size="sm" onclick={addLanguage}>
              <Plus size={14} class="mr-1.5" />
              Add
            </Button>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          {#if languages.length === 0}
            <p class="text-sm text-muted-foreground text-center py-4">No languages added yet. Click "Add" to list languages you speak.</p>
          {/if}
          {#each languages as lang, i}
            <div class="border border-border rounded-lg p-4 relative">
              <button
                type="button"
                onclick={() => removeLanguage(i)}
                class="absolute top-3 right-3 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
                aria-label="Remove language"
              >
                <X size={16} />
              </button>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pr-8">
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Language</label>
                  <Input type="text" bind:value={languages[i].language} placeholder="English" />
                </div>
                <div class="space-y-1.5">
                  <label class="text-xs font-medium text-muted-foreground">Proficiency</label>
                  <select
                    bind:value={languages[i].proficiency}
                    class="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 transition-colors"
                  >
                    <option value="native">Native</option>
                    <option value="fluent">Fluent</option>
                    <option value="advanced">Advanced</option>
                    <option value="intermediate">Intermediate</option>
                    <option value="beginner">Beginner</option>
                  </select>
                </div>
              </div>
            </div>
          {/each}
        </CardContent>
      </Card>

      <!-- Achievements -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-yellow-500/10 flex items-center justify-center">
              <Trophy size={16} class="text-yellow-500" />
            </div>
            <div>
              <CardTitle class="text-base">Achievements</CardTitle>
              <CardDescription class="text-xs">Awards, recognitions, and notable accomplishments</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <textarea
              bind:value={achievements}
              rows="4"
              placeholder="List your key achievements, awards, and recognitions..."
              class={textareaClass}
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Save Button -->
      <div class="flex justify-end pt-2 animate-fade-in-up">
        <Button type="submit" disabled={saving}>
          {#if saving}
            <Loader2 size={16} class="mr-2 animate-spin" />
            Saving...
          {:else}
            <Save size={16} class="mr-2" />
            Save Profile
          {/if}
        </Button>
      </div>
    </form>
  {/if}
</div>
