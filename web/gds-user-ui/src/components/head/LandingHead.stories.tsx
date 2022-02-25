import React from "react";
import { Story } from "@storybook/react";
import LandingHead from "./LandingHead";

interface LandingHeadProps {}

export default {
  title: "Landing",
  component: LandingHead,
};

export const Head: Story<LandingHeadProps> = ({ ...props }) => (
  <LandingHead {...props} />
);

Head.bind({});
