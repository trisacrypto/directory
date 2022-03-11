import { Meta, Story } from "@storybook/react";

import LegalPersonReview from "./LegalPersonReview";

export default {
  title: "components/LegalPersonReview",
  component: LegalPersonReview,
} as Meta;

const Template: Story = (args) => <LegalPersonReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
