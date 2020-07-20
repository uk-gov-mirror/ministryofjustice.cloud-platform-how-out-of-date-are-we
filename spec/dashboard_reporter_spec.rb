require "spec_helper"

describe DashboardReporter do
  let(:data) { {
    "updated_at" =>"2020-07-20 09:27:44",
    "data" => {
      "action_items" => {
        "documentation" => 1,
        "helm_whatup" => 1,
        "repositories" => 3,
        "terraform_modules" => 0
      },
      "action_required" => action_required
    }
  } }

  let(:formatted_message) { %(
How out of date are we - action required:
```
documentation:     1
helm_whatup:       1
repositories:      3
terraform_modules: 0
```
                            ).strip }

  let(:dashboard_url) { "" }

  subject(:dr) { described_class.new(dashboard_url) }

  before do
    allow(dr).to receive(:data).and_return(data)
  end

  context "when there are no open todo items" do
    let(:action_required) { false }

    it "returns empty string" do
      expect(dr.report).to eq("")
    end
  end

  context "when there are open todo items" do
    let(:action_required) { true }

    it "formats the report for posting to slack" do
      expect(dr.report).to eq(formatted_message)
    end
  end

  context "when data is incorrectly structured" do
    let(:data) { { "foo" => "bar" } }

    it "raises an error" do
      expect{
        dr.report
      }.to raise_error(KeyError)
    end
  end
end
