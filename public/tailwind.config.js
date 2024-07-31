tailwind.config = {
  theme: {
    extend: {
      colors: {
        link: "hsl(120, 20%, 25%)",
      },
      borderColor: {
        DEFAULT: "rgba(0, 0, 0, 0.15)",
      },
    },
    container: {
      center: true,
      screens: {
        sm: "100%",
        md: "100%",
        lg: "100%",
        xl: "968px",
        "2xl": "968px",
      },
      padding: {
        DEFAULT: "1rem",
        xl: "2rem",
      },
    },
  },
};
