tailwind.config = {
  theme: {
    extend: {
      colors: {
        ghost: "hsl(120, 20%, 88%)",
        info: "hsl(120, 20%, 88%)",
        warning: "hsl(9, 69%, 88%)",
        danger: "hsl(0, 100%, 90%)",
        success: "hsl(133, 75%, 88%)",
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
        xl: "1080px",
        "2xl": "1080px",
      },
      padding: {
        DEFAULT: "1rem",
        xl: "2rem",
      },
    },
  },
};
