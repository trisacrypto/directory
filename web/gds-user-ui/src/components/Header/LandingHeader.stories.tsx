import React from "react";
import { Story } from "@storybook/react";

import  LandingHeader  from './LandingHeader';

interface LandingHeaderProps {

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

Header.bind({});