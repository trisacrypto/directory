import { Meta, Story } from "@storybook/react";

import BasicDetailsSection from "./BasicDetailsSection";

export default {
  title: "components/BasicDetailsSection",
  component: BasicDetailsSection,
} as Meta;

const Template: Story = (args) => <BasicDetailsSection {...args} />;

export const Default = Template.bind({});
Default.args = {};
