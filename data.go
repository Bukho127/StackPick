package main

// Library represents a single installable package (or group of npm packages
// that always travel together, e.g. Tailwind + postcss + autoprefixer).
type Library struct {
	Name        string
	Packages    []string
	Description string
}

// Category groups related libraries together in the checklist UI.
type Category struct {
	Name string
	Libs []Library
}

// FrontendCategories is the catalog used when the user selects frontend.
var FrontendCategories = []Category{
	{
		Name: "Routing",
		Libs: []Library{
			{"React Router", []string{"react-router-dom"}, "Declarative client-side routing"},
		},
	},
	{
		Name: "State Management",
		Libs: []Library{
			{"Zustand", []string{"zustand"}, "Small, fast, unopinionated state store"},
			{"Redux Toolkit", []string{"@reduxjs/toolkit", "react-redux"}, "Official, batteries-included Redux"},
			{"Jotai", []string{"jotai"}, "Primitive and flexible atomic state"},
		},
	},
	{
		Name: "Data Fetching",
		Libs: []Library{
			{"TanStack Query", []string{"@tanstack/react-query"}, "Async state, caching & sync"},
			{"SWR", []string{"swr"}, "React hooks for remote data fetching"},
			{"Axios", []string{"axios"}, "Promise-based HTTP client"},
		},
	},
	{
		Name: "Styling",
		Libs: []Library{
			{"Tailwind CSS", []string{"tailwindcss", "postcss", "autoprefixer"}, "Utility-first CSS framework"},
			{"styled-components", []string{"styled-components"}, "CSS-in-JS with tagged templates"},
			{"Emotion", []string{"@emotion/react", "@emotion/styled"}, "Performant, flexible CSS-in-JS"},
		},
	},
	{
		Name: "Forms & Validation",
		Libs: []Library{
			{"React Hook Form", []string{"react-hook-form"}, "Performant, minimal-re-render forms"},
			{"Formik", []string{"formik"}, "Build forms without the tears"},
			{"Zod", []string{"zod"}, "TypeScript-first schema validation"},
		},
	},
	{
		Name: "Animation",
		Libs: []Library{
			{"Framer Motion", []string{"framer-motion"}, "Production-ready motion library"},
		},
	},
	{
		Name: "UI & Icons",
		Libs: []Library{
			{"Lucide Icons", []string{"lucide-react"}, "Clean, consistent icon set"},
			{"Radix UI", []string{"radix-ui"}, "Unstyled, accessible UI primitives"},
		},
	},
	{
		Name: "Utilities",
		Libs: []Library{
			{"Lodash", []string{"lodash"}, "General purpose utility functions"},
			{"date-fns", []string{"date-fns"}, "Modern date utility library"},
			{"clsx", []string{"clsx"}, "Tiny utility for conditional classNames"},
		},
	},
}

// BackendCategories is the catalog used when the user selects backend.
var BackendCategories = []Category{
	{
		Name: "Servers",
		Libs: []Library{
			{"Express", []string{"express"}, "Minimal and flexible Node.js web framework"},
			{"Fastify", []string{"fastify"}, "Fast and low-overhead web framework"},
		},
	},
	{
		Name: "ORMs",
		Libs: []Library{
			{"Prisma", []string{"prisma", "@prisma/client"}, "Modern database toolkit and ORM"},
			{"Sequelize", []string{"sequelize"}, "Feature-rich SQL ORM for Node.js"},
		},
	},
	{
		Name: "Authentication",
		Libs: []Library{
			{"Passport", []string{"passport"}, "Authentication middleware for Node.js"},
			{"JWT", []string{"jsonwebtoken"}, "JSON Web Token signing and verification"},
		},
	},
	{
		Name: "Validation",
		Libs: []Library{
			{"Joi", []string{"joi"}, "Schema description and data validation"},
			{"Zod", []string{"zod"}, "TypeScript-first schema validation"},
		},
	},
	{
		Name: "Testing",
		Libs: []Library{
			{"Jest", []string{"jest"}, "Delightful JavaScript testing framework"},
			{"Vitest", []string{"vitest"}, "Fast Vite-native test runner"},
		},
	},
}

// Categories is kept as a compatibility alias for the initial catalog.
var Categories = FrontendCategories
