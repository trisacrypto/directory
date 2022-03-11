import { Meta, Story } from "@storybook/react";

import TrixoReview from "./TrixoReview";

export default {
  title: "components/TrixoReview",
  component: TrixoReview,
} as Meta;

const Template: Story = (args) => <TrixoReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
