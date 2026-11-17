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

type Testimonial struct{ Quote, Name, Course string }

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
	MapEmbedURL: "", // TODO
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
		{"School Clubs", ""},
		{"Term Programme", ""},
		{"Holiday Boot Camp", ""},
		{"School Digital Academy", ""},
		{"Online", ""},
	},
	WhyUs: []string{
		"100% Practical",
		"Portfolio-First",
		"Certificate on completion",
		"Project-based learning",
		"AI-integrated skills",
	},
	Socials: []Social{
		{"Instagram", ""}, // TODO
		{"Facebook", ""},  // TODO
		{"TikTok", ""},    // TODO
		{"LinkedIn", ""},  // TODO
		{"YouTube", ""},   // TODO
	},
	Testimonials: []Testimonial{}, // real quotes only, with permission
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