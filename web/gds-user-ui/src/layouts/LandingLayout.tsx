import React from "react";
import { Flex, FlexProps } from "@chakra-ui/react";
import LandingHeader from "components/Header/LandingHeader";
import LandingHead from "components/Head/LandingHead";
import AboutTrisaSection from "components/Section/AboutUs";
import JoinUsSection from "components/Section/JoinUs";
import LandingFooter from "components/Footer/LandingFooter";

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
