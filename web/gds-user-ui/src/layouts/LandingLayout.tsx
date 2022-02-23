import React from "react";
import { Flex, FlexProps } from "@chakra-ui/react";
import LandingHeader from "../components/header/LandingHeader"
import LandingFooter from '../components/footer/LandingFooter';



export default function LandingLayout(props : FlexProps) : JSX.Element {
  return (
    <Flex
      direction="column"
      align="center"
      maxW={{ xl: "1200px" }}
      m="0 auto"
      {...props}
    >
      <LandingHeader />
          
          {props.children}
          
       <LandingFooter />
      </Flex>
      
      
  );
}