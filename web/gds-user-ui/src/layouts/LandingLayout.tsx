import React from "react";
import { Flex, FlexProps } from "@chakra-ui/react";
import Header from "../components/Header/LandingHeader";




export default function LandingLayout(props : FlexProps) {
  return (
    <Flex
      direction="column"
      align="center"
      maxW={{ xl: "1200px" }}
      m="0 auto"
      {...props}
    >
      <Header />
      {props.children}
    </Flex>
  );
}