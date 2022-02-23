import React from "react";
import { Story } from "@storybook/react";

import  LandingFooter  from './LandingFooter';

interface LandingFooterProps {

}

export default {
  title: "Landing",
  component: LandingFooter,
};

export const Footer: Story<LandingFooterProps> = ({ ...props }) => (
  <LandingFooter
    {...props}
  />
);

Footer.bind({});