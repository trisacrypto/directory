import { extendTheme } from "@chakra-ui/react";

const colors = {
    transparent: "transparent",
    black: "#000",
    white: "#fff",
    gray: {
        50: "#f7fafc",
        900: "#171923",
    },
    system: {
        blue: "#23A7E0",
        cyan: "#1BCE9F",
        orange: "#FF7A59",
        white: "#E3EBEF",
        gray: "#5B5858",
        green: "#0A864F",
    }
};

const customTheme = extendTheme({ colors });

export default customTheme;