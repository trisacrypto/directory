import { Meta, Story } from "@storybook/react";
import TrisaVerifiedLogo from ".";

export default {
  title: "components/TrisaVerifiedLogo",
  component: TrisaVerifiedLogo,
} as Meta;

const Template: Story = (args) => <TrisaVerifiedLogo {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
