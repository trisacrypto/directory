import React from "react";
import { Story } from "@storybook/react";
import IntegrateAndComply from "./IntegrateAndComply";

interface IntegrateAndComplyProps {}

export default {
  title: "Components/IntegrateAndComply",
  component: IntegrateAndComply,
};

export const Default: Story<IntegrateAndComplyProps> = ({ ...props }) => (
  <IntegrateAndComply {...props} />
);

Default.bind({});
