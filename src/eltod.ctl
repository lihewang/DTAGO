TITLE                                   ELToD 4.0 SERPM Y2045
                                        
PROJECT_DIRECTORY                       
DEFAULT_FILE_FORMAT                     COMMA_DELIMITED
PAGE_LENGTH                             99999
MODEL_START_TIME                        0.0
MODEL_END_TIME                          24.0
MODEL_TIME_INCREMENT                    15 minutes
                                        
LINK_FILE                               I:\ELToD4\Base\SF\Y2045\LINK.CSV
LINK_FORMAT                             COMMA_DELIMITED
EXPRESS_FACILITY_TYPES                  96,99,97,98
EXPRESS_LINK_ENTRY_TYPE                 97
EXPRESS_LINK_EXIT_TYPE                  98                             
LINK_TOD_CLOSER_FILE                    I:\ELToD4\Input\SF\2045\Link_TOD_Closer.csv
LINK_TOD_CLOSER_FORMAT                  COMMA_DELIMITED
                                        
TURN_PROHIBITION_FILE                   I:\ELToD4\Input\SF\2045\Turn_Prohibit.csv
TURN_PROHIBITION_FORMAT                 COMMA_DELIMITED
                                        
NODE_FILE                               I:\ELToD4\Base\SF\Y2045\NODE.CSV
NODE_FORMAT                             COMMA_DELIMITED
EXPRESS_ENTRY_TYPES                     90,92
EXPRESS_EXIT_TYPES                      91,92
GENERAL_JOIN_TYPES                      93
ZONE_NODE_TYPE                          99
                                        
TRIP_FILE                               I:\ELToD4\GO\TT_FINAL.csv
TRIP_FORMAT                             DBASE
MINIMUM_TRIP_SPLIT                      0.01
STORE_TRIPS_IN_MEMORY                   FALSE
                                        
TOLL_FILE                               I:\ELToD4\Input\SF\2045\Toll_Link.csv
TOLL_FORMAT                             COMMA_DELIMITED
                                        
TOD_TOLL_FILE                           I:\ELToD4\Input\SF\2045\TOD_Toll_Data.csv
TOD_TOLL_FORMAT                         COMMA_DELIMITED
                                        
TOLL_CONSTANT_FILE                      I:\ELToD4\Input\SF\2045\Toll_Constants.csv
TOLL_CONSTANT_FORMAT                    COMMA_DELIMITED
                                        
NUMBER_OF_THREADS                       2
MAXIMUM_ITERATIONS                      30
TRAVEL_TIME_CONVERGENCE                 0.02
EXPRESS_TOLL_CONVERGENCE                0.02
IMPEDANCE_CONVERGENCE                   0.02
MINIMUM_SPEED                           5
                                        
DISTANCE_VALUE                          0.0
TIME_VALUE                              1
COST_VALUE                              3.73,3.73,3.73,2.32,1.09,1.68                     
MODE_COST_FACTORS                       1.0, 1.0, 1.0, 1.0, 1.0, 3.1
MODE_PCE_FACTORS                        1.0, 1.0, 1.0, 1.0, 1.0, 1.5
                                        
TOLL_POLICY_CODES                       1,2,3,4
MINIMUM_TOLL                            0.5,0.5,0.5,0.5
MAXIMUM_TOLL                            10.50,10.50,10.50,1.00
MAXIMUM_VC_RATIO                        5.0,5.0,5.0,5.0
VC_RATIO_OFFSET                         0.1,0.3,0.05,0.0
TOLL_EXPONENT                           6.5,3.5,4.5,8.5
MAXIMUM_TOLL_CHANGE                     10,10,10,10
                                        
MODEL_TIME_FACTOR                       -0.0382, -0.08836, -0.13527, -0.17698,-0.20036,-0.17105
MODEL_TOLL_FACTOR                       -0.54999, -0.54783, -0.50519, -0.41070,-0.21818,-0.28804
MODEL_RELIABILITY_RATIO                 3
MODEL_RELIABILITY_TIME                  0.2
MODEL_RELIABILITY_DISTANCE              0.1
MODEL_PERCEIVED_TIME                    13.67
MODEL_PERCEIVED_MID_VC                  0.693
MODEL_PERCEIVED_MAX_VC                  1.65
MODEL_EXPRESS_WEIGHT                    1.28
MODEL_SCALE_LENGTH                      7.2
MODEL_SCALE_ALPHA                       0
MODEL_MAX_CIRCUITY                      1.5
                                        
PATH_PERCEIVED_TIME                     13.67
PATH_PERCEIVED_MID_VC                   0.693
PATH_PERCEIVED_MAX_VC                   1.65
                                        
SMOOTH_GROUP_SIZE                       3
PERCENT_MOVED_FORWARD                   20
PERCENT_MOVED_BACKWARD                  20
SMOOTHING_ITERATIONS                    6
CIRCULAR_GROUP_FLAG                     TRUE
DAILY_WRAP_FLAG                         TRUE
ITERATION_VOLUME_FLAG                   FALSE
DUMP_PARAMETER_DATA                     TRUE
IMPEDANCE_SORT_METHOD                   TRUE
                                        
NEW_VOLUME_FILE                         I:\ELToD4\GO\Volume.csv
NEW_VOLUME_FORMAT                       COMMA_DELIMITED
NEW_MODEL_DATA_FILE                     I:\ELToD4\GO\Choice_Model_Log_File.csv
NEW_MODEL_DATA_FORMAT                   COMMA_DELIMITED
SELECT_MODEL_PERIODS                    1
SELECT_MODEL_ITERATIONS                 1
SELECT_MODEL_MODES                      AUTOVOT1,AUTOVOT5
SELECT_MODEL_NODES                      100064
                                        
NEW_PERIOD_GAP_FILE                     I:\ELToD4\Base\SF\Y2045\period_gap_file.csv
NEW_PERIOD_GAP_FORMAT                   COMMA_DELIMITED
NEW_TIME_GAP_FILE                       I:\ELToD4\Base\SF\Y2045\time_gap_file.csv
NEW_TOLL_GAP_FILE                       I:\ELToD4\Base\SF\Y2045\toll_gap_file.csv
NEW_IMPEDANCE_GAP_FILE                  I:\ELToD4\Base\SF\Y2045\impedance_gap_file.csv
NEW_EXPRESS_TOLL_FILE                   I:\ELToD4\Base\SF\Y2045\express_toll_file.csv
NEW_EXPRESS_TOLL_FORMAT                 COMMA_DELIMITED
ELTOD_REPORT_1                          TIME_GAP_REPORT
ELTOD_REPORT_2                          TOLL_GAP_REPORT
ELTOD_REPORT_3                          IMPEDANCE_GAP_REPORT
ELTOD_REPORT_4                          CHOICE_DISTRIBUTION
