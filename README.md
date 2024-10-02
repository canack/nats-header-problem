# NATS Go Client Header Bug Demonstration

This repository demonstrates a bug in the NATS Go client where setting headers on messages that initially lack them causes a panic.

## Overview

In NATS, messages can have headers for storing metadata. The issue occurs when a consumer attempts to add or modify headers on a message that originally does not have any headers, leading to an application crash.

## Problem Description

### Steps to Reproduce

1. **Producer**: Sends messages to a NATS subject.
    - The first three messages include a "Publish-Count" header.
    - Subsequent messages do not include any headers.

2. **Consumer**: Receives messages and attempts to set a new header "No-Panic":
    - If the message already has headers, the operation succeeds.
    - If the message has no headers, the operation causes a panic due to the lack of header initialization.

### Illustration

```plaintext  
  Scenario 1: Message with Initial Headers  
  +-----------------+                       +--------------------+  
  |  Producer       |    Send Message       |  Consumer          |  
  |  (with headers) |---------------------->|  Reads Headers     |  
  +-----------------+                       |  Set/Add Header    |  
                                            |  (success)         |  
                                            +--------------------+  

  Scenario 2: Message without Initial Headers  
  +-----------------+                       +--------------------+  
  |  Producer       |    Send Message       |  Consumer          |  
  | (no headers)    |---------------------->|  Attempt Set/Add   |  
  +-----------------+                        |  Header (panic!)  |  
                                            +--------------------+  
```

## Solution

To fix this issue, ensure that the consumer checks if the message contains headers before setting new ones. If not, initialize the header map to prevent a panic and maintain application stability.
