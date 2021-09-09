
UPDATE_COMPLETE
CREATE_FAILED
CREATE_COMPLETE

subcommands:
  - apply
      - var-file
      - -auto-approve
  - destroy
  - validate
  - plan (X)
  - show -> Show the current state or a saved plan
  - list-exports
  - fmt (?)
  - version

- plan 
cfnctl has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

var.image_id
  Enter a value: