import React from "react";
import { Story } from "@storybook/react";

import { MobileNavBar } from './MobileNav';


interface MobileNavBarProps{


}
export default {
  title: "NavBar/MobileNav",
  component: MobileNavBar,
};

export const Landing: Story<MobileNavBarProps> = ({ ...props }) => (
  <MobileNavBar
    {...props}
  />
);

Landing.bind({});