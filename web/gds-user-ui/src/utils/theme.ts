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
        link: "#23A7E0",
        btng: "#1BCE9F",
        btno: "#FF7A59",
        wht: "#E3EBEF",
        grey: "#5B5858",
        gr: "#0A864F",
    }
};

const customTheme = extendTheme({ colors });

export default customTheme;