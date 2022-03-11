import { Meta, Story } from "@storybook/react";

import ContactsReview from "./ContactsReview";

export default {
  title: "components/ContactsReview",
  component: ContactsReview,
} as Meta;

const Template: Story = (args) => <ContactsReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
