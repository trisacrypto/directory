import React from "react";
import { Story } from "@storybook/react";
import LandingHead from "./LandingHead";

interface LandingHeadProps {}

export default {
  title: "components/LandingHead",
  component: LandingHead,
};

export const Default: Story<LandingHeadProps> = ({ ...props }) => (
  <LandingHead {...props} />
);

Default.bind({});
