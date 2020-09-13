
0.9.17
=============
2020-09-13

* Allowed new record type in the DNS resource (f1524b4b)
* Updated the civogo lib to the v0.2.17 (38a311e3)
* Merge pull request #22 from fgrehm/patch-1 (389e6c06)
* Update the readme and clean go.mod (1171aedb)
* Update the Change log (26522e57)
* Update the release action (e36ac2a4)
* Update the release to use Go 1.14 (5cfb6289)
* Update the goreleaser config (6e381d00)
* Update Change log (932440f7)

0.9.16
=============
2020-09-02

* Update the release action (e36ac2a4)
* Update the release to use Go 1.14 (5cfb6289)
* Update the goreleaser config (6e381d00)
* Update Change log (932440f7)
* Added GPG to sing the binary (942f59cd)
* Update docs to the new format (a3ee7776)

0.9.15
=============
2020-08-17

* Fixed bug in the kubernetes cluster resource (bee9eb41)
* Update the README, remove Progress resource section (4a792b9a)
* Updated the Change log (32c5e71f)
* Fixed error in the description of the instance resource (9314aa97)

0.9.14
=============
2020-08-03

* Fixed error in the description of the instance resource (9314aa97)
* Updated the Change log (6db9f524)
* Fixed error in the instance resource (1df99212)

0.9.13
=============
2020-07-17

* Changed the direction in the firewall rule from inbound to ingress (cd24e0a9)
* Update goreleaser-action to v2 (a50103ad)
* Update kubernetes cluster data source with the new fields (94bd2d64)
* Update change log (eb0123f1)

0.9.12
=============
2020-07-07

* Added CPU, RAM and SSD fields to Instance and Kubernetes module (d5eaef11)
* Update change log (6e56ab4b)

0.9.11
=============
2020-07-06

* Added kubernetes cluster data source (bbb18219)
* Added the make test to the github actions (92073abd)
* Update the change logs (0202a445)

0.9.10
=============
2020-07-06

* Fixed error in the kubernetes cluster, the master's ip was not set (8ab90738)
* Added the change logs (d049062e)

0.9.9
=============
2020-07-06

* Fixed error in the kubernetes resource (9bd1bc71)
* Update the documentation of the kubernetes cluster (4e5b922d)
* Added master_ip to kubernetes cluster (5c1ad36e)
* Added data source kubernetes versions test (9d480c9c)
* Added new test (4f3a1a1a)
* Added new test (834a4e30)
* Fixed error in the doc (8b75b833)
* Added data source snapshot test (12c41ce8)
* Added data source Template test (efbfce8f)
* Added new data source test (439c0d42)
* Added new test to the provider (26017174)
* Added new test (91509d2d)
* Added the Template test (102a5b22)
* Fixed error in the doc again (ca85af3a)
* Fixed error in the doc of the template resource (5c5c95db)
* Fixed kubernetes cluster (68cca033)
* Added new test (75e69a84)
* Fixed error in provider (0c87d356)
* Fixed error in the provider (b73725b8)
* Added the SSH Key test and the Snapshot test (61987e1e)

0.9.8
=============
2020-06-24

* Upgrade the version of civogo lib (a3077b4e)
* Fixed error ins some test (25a1a6b2)
* Merge pull request #18 from AugustasV/patch-1 (1e3c6658)

0.9.7
=============
2020-06-22

* Fixed error in the instances creation (bd3664d9)
* test: Add new test (961f2831)
* Added new test (3a237d79)
* Add DNS domain record test (334cdb57)
* Added the first test to the project (400524a0)
* Update documetation (2127d822)

0.9.6
=============
2020-06-07

* Update the documentation (dec82903)
* Change the provider, now the token is not required (a14c2e59)
* Added the option to recreate resource if not found (256b2e36)
* Fix error in k8s cluster after deletion (a239402e)
* Fixed error in all this files (0bb57878)
* Fix some error found in the code (05ec3ecb)

0.9.5
=============
2020-05-31

* Added a error handler to the provider in the case of the token not was found (a456dad3)
* Added default value to size and initial_user in the instance resource (4dfc12da)
* - Update the documentation for data source instance_size and template (549da91b)

0.9.4
=============
2020-05-21

* - Added the new data source snapshot (96f9e41c)
* Added the new data source ssh key to get one ssh from the civo cloud (e32df07e)
* Added the new data source loadbalancer to get one loadbalancer from the civo cloud (255aaf46)
* feat(DataSource): Added new data source (c97d499b)
* feat(Doc): Fix error in doc (69cb9293)
* feat(DataSource): Added new data source (cf116faa)
* feat(DataSource): Added new data source (aeb39320)
* feat(DataSource): Added new data source (10629d6d)
* feat(DataSource): Added new data source (559594ae)
* feat(DataSource): Added new data source (40a16547)
* Change in the readme (465ad8bd)

0.9.3
=============
2020-05-05

* Added support to add script (807f086e)

0.9.2
=============
2020-04-24

* improved the data source filter (197425b7)

0.9.1
=============
2020-04-20

* Merge pull request #9 from civo/dev (2ba9c9ef)
* Merge pull request #8 from civo/dev (2a8e4c75)

0.9.0
=============
2020-04-18

* Pre release (4172f90a)
* Merge pull request #6 from civo/dev (4568fde4)
* Add status badge for GitHub actions (dd8e95d7)
* Remove Travis integration (d2dba75c)
* Use GitHub actions (862f18a8)
* Create CODEOWNERS (95a61bc5)
* Merge pull request #4 from jdbohrman/patch-1 (5781356a)
* Merge pull request #3 from alejandrojnm/dev (60387abc)
* Merge pull request #2 from alejandrojnm/dev (e392c9dc)
* Update README to add status (d9e6c92f)
* Remove website test until we're on there (7df84a35)
* Add vendored sources (94afb93d)
* Renamed Makefile (93d44b42)
* Add travis status to title (4340f345)
* Template files added (e8e23f4b)
* Initial readme update (7845305f)
* Initial commit (718385ef)
* - Add doc.md to know how to use the provider (2ab8c472)
* - Remove tags from the instance creation for now (935ec412)
* - Remove validator from resource_instance.go to utils.go (47fc803f)
* - Add validate to the instance reverse_dns (no white spaces) (89aa8636)
* - Add the instances delete function (14822c95)
* - Add the instances red function (db611b65)
* First commit for the terraform provider for civo (e018fda8)
* - Update .travis.yml (4382fc45)
* - Rename provider folder to civo (237f68b5)
* Fix doc.md typo? (ea7666f5)
* fix: Remove a temporal file (fd394138)
* feat: Added one feature to the provider (e2e5c6f8)
* feat: Added one feature to the snapshot resource (52bf3e07)
* feat: Added the option to import existing infrastructure (e4dc0112)
* feat: Add other data resource, k8s version, instances size (77c3d358)
* feat: Add other resource, kubernet, snapshot (08929849)
* feat: Add templates module (79ff73a9)
* feat: Add ssh resource (04c84fe8)
* feat: Add firewall and loadbalancer (834be2e1)
* - Add civo_dns_domain_name and civo_dns_domain_record (73ebfd3d)
* - Fix error in tags (7f7d75df)
* - Add tags to the instance (e5e45e64)
* Merge branch 'master' into dev (62c4abe3)
* - Update civo go api (2552d950)
* - Update civo go api (974af8b7)
* - Update go 1.13 (95287034)
* Merge remote-tracking branch 'origin/master' (0e81005a)
* - Update go 1.13 (310145cc)
* Merge pull request #1 from civo/master (1c8c2b6b)
* feat: Add goreleaser conf (10be43bc)


