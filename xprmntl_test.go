package xprmntl

import (
	"errors"
	"testing"
	"os"
	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestNew(t * testing.T) {
	os.Setenv("FEATURE_URL", "http://example.com/");
	os.Setenv("FEATURE_DEVKEY", "exp-dev-key-here");
	os.Setenv("FEATURE_DEVKEY_SHARED", "exp-dev-key-here");

	// Setup Testing Framework
	g := Goblin(t);
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) });

	g.Describe(".config interface:", func() {
			g.Describe("Given no config object,", func() {
					g.It("should reject with an Error", func() {
							cli, err := New(nil);
							Expect(cli).Should(BeNil());
							Expect(err).Should(Equal(errors.New("XPRMNTL: New(): Cannot register experiments without a config. Please see docs.")));
						});
				});
			g.Describe("Given no experiments array", func() {
					g.It("should reject with an Error", func() {
							config := new(Config);
							cli, err := New(config);
							Expect(cli).Should(BeNil());
							Expect(err).Should(Equal(errors.New("XPRMNTL: New(): Cannot register experiments without `experiments`. Please see the docs.")));
						})
				});
			g.Describe("Given no experiments in array", func() {
					g.It("should reject with an Error", func() {
							config := Config{Experiments: []Experiment{}};
							cli, err := New(&config);
							Expect(cli).Should(BeNil());
							Expect(err).Should(Equal(errors.New("XPRMNTL: New(): Cannot register experiments without `experiments`. Please see the docs.")));
						});
				});
			g.Describe("Given no optional Configurations", func() {
					g.It("should fall back to ENV", func() {
							config := Config {
							Experiments: []Experiment {
								Experiment{Name: "testEx"},
							},
						};
							cli, err := New(&config);
							Expect(err).Should(BeNil());
							Expect(*cli.DevKey).Should(Equal(os.Getenv("FEATURE_DEVKEY")));
							Expect(*cli.FeatureURL).Should(Equal(os.Getenv("FEATURE_URL")));
							Expect(cli.Shared.DevKey).Should(Equal(os.Getenv("FEATURE_DEVKEY_SHARED")));
						});
				});
			g.Describe("Given a passed devKey / featureUrl,", func() {
					g.It("should not fall back to ENV", func() {
							config := Config {
							DevKey: "new-dev-key-here",
							FeatureURL: "https://example.com/",
							Experiments: []Experiment {
								Experiment{Name: "testEx"},
							},
						};
							cli, err := New(&config);
							Expect(err).Should(BeNil());
							Expect(*cli.DevKey).Should(Equal("new-dev-key-here"));
							Expect(*cli.FeatureURL).Should(Equal("https://example.com/"));
						});
				});
			g.Describe("Given no timeout,", func() {
					g.It("should set timeout to be 5s", func() {
							config := Config {
							DevKey: "new-dev-key-here",
							FeatureURL: "https://example.com/",
							Experiments: []Experiment {
							Experiment{Name: "testEx"},
						},
						};
							cli, err := New(&config);
							Expect(err).Should(BeNil());
							Expect(cli.Timeout).Should(Equal(5000));
						});
				});
		});
};
