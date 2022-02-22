import React from "react";
import { Story } from "@storybook/react";

import  LandingHeader  from './LandingHeader';

interface LandingHeaderProps {

}

export default {
  title: "Landing/Header",
  component: LandingHeader,
};

export const Landing: Story<LandingHeaderProps> = ({ ...props }) => (
  <LandingHeader
    {...props}
  />
);

Landing.bind({});