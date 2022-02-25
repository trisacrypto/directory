import React from "react";
import { Story } from "@storybook/react";
import LandingHeader from "./LandingHeader";


interface LandingHeaderProps {
  title: string;
  description?: string;
}

export default {
  title: "Landing",
  component: LandingHeader,
};

export const Header: Story<LandingHeaderProps> = ({ ...props }) => (
  <LandingHeader
    {...props}
  />
);

Header.bind({
  title: "Landing",
  descirption:"description"
});