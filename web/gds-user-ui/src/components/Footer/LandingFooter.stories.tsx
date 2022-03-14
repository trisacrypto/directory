import React from "react";
import { Story } from "@storybook/react";
import LandingFooter from "./LandingFooter";
interface LandingFooterProps {}

export default {
  title: "components/LandingFooter",
  component: LandingFooter,
};

export const Default: Story<LandingFooterProps> = ({ ...props }) => (
  <LandingFooter {...props} />
);

Default.bind({});
