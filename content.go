package main

import "net/url"

type Course struct {
	Num   int
	Title string
}

// No price field on purpose: every course says "Contact for Price".
func (c Course) EnquireURL() string {
	return "/contact?interest=" + url.QueryEscape(c.Title)
}

type Social struct{ Name, URL string }

type Testimonial struct {
	ID     int
	Quote  string
	Name   string
	Email  string // <-- ADDED
	Course string
	Rating int
}

type Option struct{ Title, Blurb string }

type SiteInfo struct {
	Company, Tagline, Slogan string
	Who, Vision, Mission     string
	Address, Phone, Email    string
	WhatsApp                 string // digits only, international format
	MapEmbedURL              string // paste the src="" from Google Maps > Share > Embed a map
	PriceLabel               string
	Courses                  []Course
	Delivery                 []Option
	WhyUs                    []string
	Socials                  []Social
	Testimonials             []Testimonial
}

var site = SiteInfo{
	Company: "NexGen Guru Consulting",
	Tagline: "Empowering the Next Generation Through AI, Creativity & Technology",
	Slogan:  "Learn. Create. Innovate. Lead.",
	Who: "NexGen Guru Consulting is a Digital Education and Creative Technology company. " +
		"We prepare young people for opportunities created by AI, digital media and the creative economy " +
		"through practical, project-based training - not just theory.",
	Vision:      "To empower the next generation of African innovators, creators and tech professionals.",
	Mission:     "To equip students with practical AI, digital and creative skills through project-driven learning.",
	Address:     "No. 18 Atomkpera Road, Ankpa Quarters, Opp. Kismet Hotel, Makurdi",
	Phone:       "0806 596 7978",
	Email:       "jameskungwa@gmail.com",
	WhatsApp:    "2348065967978",
	MapEmbedURL: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d15758.019612835146!2d8.513233900000001!3d7.7325255!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x10516681a12d2e8f%3A0x6b2c6a6d3e6b8c0!2sMakurdi%2C%20Benue%20State!5e0!3m2!1sen!2sng!4v1700000000000!5m2!1sen!2sng",
	PriceLabel:  "Contact for Price",

	Courses: []Course{
		{1, "AI-Integrated Content Creation"},
		{2, "Graphic Design & Digital Branding"},
		{3, "Audio Editing & Sound Engineering"},
		{4, "Digital & AI-Supported Music Production"},
		{5, "Video Editing & Motion Graphics"},
		{6, "Professional Subtitling & Captioning"},
		{7, "Digital Entrepreneurship & Freelancing"},
	},
	Delivery: []Option{
		{"School Clubs", "On-campus after-school programs."},
		{"Term Programme", "Integrated within the school term."},
		{"Holiday Boot Camp", "Intensive holiday training."},
		{"School Digital Academy", "Full digital curriculum setup."},
		{"Online", "Learn from anywhere in the world."},
	},
	WhyUs: []string{
		"100% Practical",
		"Portfolio-First",
		"Certificate on completion",
		"Project-based learning",
		"AI-integrated skills",
	},
	Socials: []Social{
		{"Instagram", "https://instagram.com/nexgenguru"},
		{"Facebook", "https://facebook.com/nexgenguru"},
		{"TikTok", "https://tiktok.com/@nexgenguru"},
		{"LinkedIn", "https://linkedin.com/company/nexgenguru"},
		{"YouTube", "https://youtube.com/@nexgenguru"},
	},
	Testimonials: []Testimonial{}, // Leave empty for now, the HTML will show a "No Testimonials Yet" state.
}

func (s SiteInfo) WhatsAppLink() string { return "https://wa.me/" + s.WhatsApp }

// Only socials that have a URL, so the page never shows dead links.
func (s SiteInfo) ActiveSocials() []Social {
	var out []Social
	for _, x := range s.Socials {
		if x.URL != "" {
			out = append(out, x)
		}
	}
	return out
}

func courseTitles() []string {
	t := make([]string, len(site.Courses))
	for i, c := range site.Courses {
		t[i] = c.Title
	}
	return t
}