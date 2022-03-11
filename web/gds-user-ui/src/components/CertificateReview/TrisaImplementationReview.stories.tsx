import { Meta, Story } from "@storybook/react";

import TrisaImplementationReview from "./TrisaImplementationReview";

export default {
  title: "components/TrisaImplementationReview",
  component: TrisaImplementationReview,
} as Meta;

const Template: Story = (args) => <TrisaImplementationReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
