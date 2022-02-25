import React from "react";
import { Story } from "@storybook/react";
import AboutTrisaSection from "./AboutUs";

interface AboutUsProps {}

export default {
  title: "Pages",
  component: AboutTrisaSection,
};

export const aboutUs: Story<AboutUsProps> = ({ ...props }) => (
  <AboutTrisaSection {...props} />
);

aboutUs.bind({});
