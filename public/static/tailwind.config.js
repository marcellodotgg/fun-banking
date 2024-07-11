tailwind.config = {
  theme: {
    extend: {
      colors: {
        ghost: "hsl(120, 20%, 88%)",
        info: "hsl(120, 20%, 82%)",
        warning: "hsl(9, 69%, 83%)",
        danger: "hsl(0, 100%, 83%)",
        success: "hsl(133, 75%, 80%)",
        link: "hsl(120, 20%, 25%)",
      },
      borderColor: {
        DEFAULT: "rgba(0, 0, 0, 0.15)",
      },
    },
    container: {
      center: true,
      padding: {
        DEFAULT: "1rem",
        xl: "2rem",
      },
    },
  },
};
