import React from "react";
import { Flex, FlexProps } from "@chakra-ui/react";
import LandingHeader from "../components/Header/LandingHeader";
import LandingFooter from '../components/Footer/LandingFooter';



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