import React from "react";
import { Story } from "@storybook/react";

import  LandingLayout  from './LandingLayout';

interface LandingLayoutProps {

}

export default {
  title: "Landing",
  component: LandingLayout,
};

export const Layout: Story<LandingLayoutProps> = ({ ...props }) => (
  <LandingLayout
    {...props}
  />
);

Layout.bind({});