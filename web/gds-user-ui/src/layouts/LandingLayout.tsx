import React from "react";
import { Flex, FlexProps } from "@chakra-ui/react";
import LandingHeader from "../components/Header/LandingHeader";
import LandingFooter from "../components/Footer/LandingFooter";
import LandingHead from "../components/head/LandingHead";
import AboutTrisaSection from "../components/section/AboutUs";
import JoinUsSection from "../components/section/JoinUs";

export default function LandingLayout(props: FlexProps): JSX.Element {
  return (
    <Flex direction="column" align="center" maxW={"100%"} m="0 auto" {...props}>
      <LandingHeader />

      {props.children}
      <LandingHead />
      <AboutTrisaSection />
      <JoinUsSection />
      <LandingFooter />
    </Flex>
  );
}
