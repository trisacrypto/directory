import { Meta, Story } from "@storybook/react";
import BasicDetailsReview from "./BasicDetailsReview";

export default {
  title: "components/BasicDetailsReview",
  component: BasicDetailsReview,
} as Meta;

const Template: Story = (args) => <BasicDetailsReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
