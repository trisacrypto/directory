import React from "react";
import { Story } from "@storybook/react";
import LandingHeader from "./LandingHeader";

interface LandingHeaderProps {
  title: string;
  description?: string;
}

export default {
  title: "components/LandingHeader",
  component: LandingHeader,
};

export const Default: Story<LandingHeaderProps> = ({ ...props }) => (
  <LandingHeader {...props} />
);

Default.bind({
  title: "Landing",
  decription: "description",
});
