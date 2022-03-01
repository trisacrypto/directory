import React from "react";
import { Story } from "@storybook/react";
import AboutTrisaSection from "./AboutUs";

interface AboutUsProps {}

export default {
  title: "Components/AboutUs",
  component: AboutTrisaSection,
};

export const standard: Story<AboutUsProps> = ({ ...props }) => (
  <AboutTrisaSection {...props} />
);

standard.bind({});


